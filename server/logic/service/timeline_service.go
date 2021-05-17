package service

import (
	"OpenFaaS-Logic/dal/db"
	"OpenFaaS-Logic/model"
	"OpenFaaS-Logic/pkg/dto"
)

type TimelineService struct{}

func NewTimelineService() *TimelineService {
	return &TimelineService{}
}

func (t *TimelineService) CreateTimeline(msg *dto.TimeLineModel) (timeline *model.Timeline, err error) {
	dalTimeline := db.NewDalTimeLine()
	timeline = &model.Timeline{
		TimelineID: msg.TimelineID,
		MessageID:  msg.MessageID,
		Type:       msg.MessageType,
		Status:     msg.Status,
		CreatedAt:  msg.CreatedAt,
		UpdatedAt:  msg.CreatedAt,
		DeletedAt:  nil,
	}
	return dalTimeline.Insert(timeline)
}

//确认消息号 ackId 为timeline_model的id
func (t *TimelineService) ACKMessage(id int64) error {
	dalProxy := db.NewDalTimeLine()
	return dalProxy.AckMessage(id)
}
