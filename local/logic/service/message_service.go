package service

import (
	"OpenFaaS-Logic/dal/db"
	"OpenFaaS-Logic/model"
	"OpenFaaS-Logic/model/responseForm"
	"OpenFaaS-Logic/pkg/dto"
	"OpenFaaS-Logic/pkg/logger"
	"OpenFaaS-Logic/st"
	"fmt"

	"sync"
	"time"
)

type MessageService struct {
	Base
}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) InitAffair(base *Base) {
	s.Base = *base
}

func (s *MessageService) CreateMessage(msgModel *dto.MessageModel) (err error) {
	dalMsg := db.NewDalMessage()
	message := &model.Message{
		DeviceID:     msgModel.DeviceID,
		FromUID:      msgModel.FromUID,
		ToID:         msgModel.ToID,
		ReceiverType: msgModel.ReceiverType,
		Type:         msgModel.MessageType,
		Status:       dto.MESSAGE_NOT_SEND,
		Content:      msgModel.Content,
		CreatedAt:    time.Now(),
	}
	err = dalMsg.CreateMessage(message)
	if err != nil {
		return
	}
	msgModel.CreatedAt = message.CreatedAt
	msgModel.MessageID = message.MessageID
	return
}

//发送一条消息
func (s *MessageService) SendMessage(msg *dto.MessageModel) error {
	timelines, err := MessageToTimeline(msg)
	if err != nil {
		return fmt.Errorf("消息发送失败 ：%v", err)
	}
	go func() {
		for _, timeline := range timelines {
			err = ProduceMessage(timeline)
			if err != nil {
				logger.Logger.Errorf("该 timeline %v id:%v 发送失败", timeline.TimelineID, timeline.ID)
			}
		}
	}()

	return nil
}

//判断用户是否在线 在线发online topic 离线发offline topic
func ProduceMessage(timelineModel *dto.TimeLineModel) error {
	var (
		err      error
		timeline *model.Timeline
	)
	srvProxy := NewTimelineService()
	timeline, err = srvProxy.CreateTimeline(timelineModel)
	if err != nil {
		st.Debug(err)
		return err
	}
	if timelineModel.TimeLineType == dto.TIMELINE_PERSON && timelineModel.ToID == timelineModel.ObjectIds {
		notifyService := NewNotifyService()
		err = notifyService.Notify(timelineModel, timeline)
		if err != nil {
			st.Debug(err)
		}
	}
	return nil
}

//message 模型转换成timeline模型
func MessageToTimeline(msg *dto.MessageModel) ([]*dto.TimeLineModel, error) {
	var err error
	wg := sync.WaitGroup{}
	timeline := make([]*dto.TimeLineModel, 0)
	lock := sync.Mutex{}
	eCh := make(chan error, 1024)
	//如果是群组，则给群组里的所有人timeline 和群组timeline 发消息
	if msg.ReceiverType == dto.RECV_PERSON { //如果是单聊 则给两个人的个人timeline 和两个人的公共timeline 发消息
		wg.Add(3)
		go getTimelineId(&wg, &lock, &timeline, msg, eCh, dto.TIMELINE_DOUBLE, msg.FromUID, msg.ToID)
		go getTimelineId(&wg, &lock, &timeline, msg, eCh, dto.TIMELINE_PERSON, msg.FromUID)
		go getTimelineId(&wg, &lock, &timeline, msg, eCh, dto.TIMELINE_PERSON, msg.ToID)

	}

	wg.Wait()
	select {
	case err = <-eCh:
		close(eCh)
		return nil, fmt.Errorf("ProduceMessage error :%v", err)
	default:
		close(eCh)
	}
	if err != nil || timeline == nil {
		return nil, fmt.Errorf("ProduceMessage error :%v", err)
	}
	return timeline, nil
}

func getTimelineId(wg *sync.WaitGroup, lock *sync.Mutex, dest *[]*dto.TimeLineModel, msg *dto.MessageModel, ch chan<- error, timeLineType int, args ...interface{}) {
	defer wg.Done()
	dalTimelineMap := db.NewDalTimelineMap()
	timelineMap, e := dalTimelineMap.GetTimeLineMap(timeLineType, args...)
	if e != nil {
		logger.Logger.Errorf("getTimelineId error: %v", e)
		ch <- e
		return
	}

	res := &dto.TimeLineModel{
		MessageModel: msg,
		FromConnID:   msg.DeviceID,
		//ToConnID:     0,
		//ID:           0,
		TimelineID: timelineMap.TimelineID,
		//Seq:          0,
		TimeLineType: timelineMap.Type,
		ObjectIds:    timelineMap.ObjectIds,
		Status:       model.TIMELINE_STATUS_NOT_READ,
		CreatedAt:    msg.CreatedAt,
	}
	if msg.FromUID == timelineMap.ObjectIds {
		res.Status = model.TIMELINE_STATUS_READ
	}
	lock.Lock()
	*dest = append(*dest, res)
	lock.Unlock()
}

//确认消息号 ackId 为timeline_model的id
func (s *MessageService) ACKMessage(id int64) error {
	dalProxy := db.NewDalTimeLine()
	return dalProxy.AckMessage(id)
}

func (s *MessageService) GetUnRead(userID string, page, pageSize int64) (form responseForm.Paginator, err error) {
	var (
		timelineMap    *model.TimelineMap
		timelines      []*model.Timeline
		simpleMessages []responseForm.SimpleMessage
	)
	dalTimelineMap := db.NewDalTimelineMap()
	timelineMap, err = dalTimelineMap.GetTimeLineMap(dto.TIMELINE_PERSON, userID)
	if err != nil {
		logger.Logger.Debug("getUnRead get timeline map err", err)
		return
	}
	dalTimeline := db.NewDalTimeLine()
	form.Page = page
	form.Size = pageSize
	form.Total = dalTimeline.GetUnReadTimelineTotalNum(timelineMap.TimelineID)
	if form.Total == 0 {
		return
	}
	timelines, err = dalTimeline.GetUnReadTimeline(timelineMap.TimelineID, page, pageSize)

	dalMessage := db.NewDalMessage()
	dalTimeline.InitAffair(s.Affair.GetTX())
	for _, timeline := range timelines {
		var message *model.Message
		message, err = dalMessage.GetMessageByMessageID(timeline.MessageID)
		if err != nil {
			return
		}
		simpleMessage := responseForm.SimpleMessage{
			MessageID:    message.MessageID,
			DeviceID:     message.DeviceID,
			FromUID:      message.FromUID,
			ToID:         message.ToID,
			ReceiverType: message.ReceiverType,
			Type:         message.Type,
			Content:      message.Content,
			Seq:          timeline.Seq,
			CreatedAt:    message.CreatedAt,
		}
		simpleMessages = append(simpleMessages, simpleMessage)
		err = dalTimeline.AckMessage(timeline.ID)
		if err != nil {
			return
		}
	}
	form.Data = simpleMessages
	return
}

func (s *MessageService) GetPersonTimelineIdByUserID(userID string) (timelineId int64, err error) {
	var timelineMap *model.TimelineMap
	dal := db.NewDalTimelineMap()
	timelineMap, err = dal.GetTimeLineMap(dto.TIMELINE_PERSON, userID)
	if err != nil {
		return
	}
	return timelineMap.TimelineID, nil
}

////获取timeline_id 对应的消息列表
//func (s *MessageService) GetTimeLine(timelineId int64, isUnread *bool) ([]*dto.TimeLineModel, error) {
//	// dalProxy := db.NewDalBind()
//
//	return nil, nil
//}

//获取timeline_id 对应的消息列表
func (s *MessageService) GetTimeLine(timelineId int64, isUnread *bool) ([]*dto.TimeLineModel, error) {

	return nil, nil
}

func (s *MessageService) GetTimeLineBySeq(timelineId, seq, page, pageSize int64) (timelines []*model.Timeline, err error) {
	dal := db.NewDalTimeLine()
	timelines, err = dal.GetTimelineBySeqWithPaginator(timelineId, seq, page, pageSize)
	return
}

func (s *MessageService) GetTimelineFormBySeq(timelineId, seq, page, pageSize int64) (form responseForm.Paginator, err error) {
	var (
		timelines      []*model.Timeline
		simpleMessages []responseForm.SimpleMessage
	)
	dalTimeline := db.NewDalTimeLine()
	form.Page = page
	form.Size = pageSize
	form.Total = dalTimeline.GetTimelineBySeqTotalNum(timelineId, seq)
	if form.Total == 0 {
		return
	}
	timelines, err = s.GetTimeLineBySeq(timelineId, seq, page, pageSize)
	if err != nil {
		logger.Logger.Debug("get timelineForm get timeline by seq error", err)
		return
	}
	dalMessage := db.NewDalMessage()
	dalTimeline.InitAffair(s.Affair.GetTX())
	for _, timeline := range timelines {
		var message *model.Message
		message, err = dalMessage.GetMessageByMessageID(timeline.MessageID)
		if err != nil {
			return
		}
		simpleMessage := responseForm.SimpleMessage{
			MessageID:    message.MessageID,
			DeviceID:     message.DeviceID,
			FromUID:      message.FromUID,
			ToID:         message.ToID,
			ReceiverType: message.ReceiverType,
			Type:         message.Type,
			Content:      message.Content,
			Seq:          timeline.Seq,
			CreatedAt:    message.CreatedAt,
		}
		simpleMessages = append(simpleMessages, simpleMessage)
		err = dalTimeline.AckMessage(timeline.ID)
		if err != nil {
			return
		}
	}
	form.Data = simpleMessages
	return
}
