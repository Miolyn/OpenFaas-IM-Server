package model

import (
	"time"
)

// Timeline struct is a row record of the timeline table in the finders_imdb database
type Timeline struct {
	ID         int64      `gorm:"primary_key;column:id;type:bigint;" json:"id"`
	TimelineID int64      `gorm:"column:timeline_id;type:bigint;" json:"timeline_id"` // timeline ID
	Seq        int64      `gorm:"column:seq;type:bigint;" json:"seq"`                 // 消息序列，一个timeline下的seq从0递增。
	MessageID  int64      `gorm:"column:message_id;type:bigint;" json:"message_id"`   // 消息ID
	Type       int        `gorm:"column:type;type:int;" json:"type"`                  // timeline类型，1.个人 2.单聊 3.群聊
	Status     int        `gorm:"column:status;type:int;" json:"status"`              // timeline状态，1.未读 2.已读
	CreatedAt  time.Time  `gorm:"column:created_at;type:datetime;" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;type:datetime;" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at;type:datetime;" json:"deleted_at"`
}

// TableName sets the insert table name for this struct type
func (t *Timeline) TableName() string {
	return "timeline"
}

// status
const (
	TIMELINE_STATUS_NOT_READ = 1 + iota
	TIMELINE_STATUS_READ
)
