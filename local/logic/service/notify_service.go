package service

import (
	"OpenFaaS-Logic/global"
	"OpenFaaS-Logic/model"
	"OpenFaaS-Logic/pkg/dto"
	"fmt"

	"github.com/imroc/req"
	"time"
)

type NotifyService struct{}

func NewNotifyService() *NotifyService {
	return &NotifyService{}
}

func (s *NotifyService) Notify(msg *dto.TimeLineModel, timeline *model.Timeline) (err error) {
	var (
		form dto.NotifyMessageModel
		resp *req.Resp
	)
	if timeline.Type != dto.TIMELINE_PERSON {
		return nil
	}
	url := fmt.Sprintf("%v/notify", global.ProxyServer)
	authHeader := req.Header{
		"Content-Type": "application/json",
	}
	form = dto.NotifyMessageModel{
		ID:           time.Now().UnixNano(),
		OpCode:       dto.OPCODE_NEW_MSG,
		FromConnID:   msg.FromConnID,
		ToConnID:     msg.ToConnID,
		FromUID:      msg.MessageModel.FromUID,
		ToID:         msg.MessageModel.ToID,
		ReceiverType: msg.ReceiverType,
		Content: dto.NotifyContent{
			Content:     msg.Content,
			MessageType: msg.MessageType,
			TimelineID:  timeline.TimelineID,
			Seq:         timeline.Seq,
			ID:          timeline.ID,
			CreatedAt:   timeline.CreatedAt,
		},
	}

	resp, err = req.Post(url, authHeader, req.BodyJSON(&form))
	if err != nil {
		return fmt.Errorf("请求connect 服务器失败:%v", err)
	}
	res := new(dto.ResponseJson)
	err = resp.ToJSON(res)
	if err != nil {
		return fmt.Errorf("resp parse error :%v", err)
	}
	if res.Code != dto.RESPONSE_OK_CODE {
		return fmt.Errorf("logic server error:%v", res.Msg)
	}
	return nil
}
