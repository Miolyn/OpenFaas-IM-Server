package model

import (
	"time"
)


type TimelineMap struct {
	TimelineID int64      `gorm:"primary_key;column:timeline_id;type:bigint;" json:"timeline_id"` // timeline 主键
	Type       int        `gorm:"column:type;type:int;" json:"type"`                              // timeline类型：1.个人 2.单聊 3.群聊
	ObjectIds  string     `gorm:"column:object_ids;type:varchar;size:255;" json:"object_ids"`     // timeline对象id: type=1 时为user_id，type=2时为user_id,user_id，type=3时为group_id
	CreatedAt  time.Time  `gorm:"column:created_at;type:datetime;" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;type:datetime;" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at;type:datetime;" json:"deleted_at"`
}

func (t *TimelineMap) TableName() string {
	return "timeline_map"
}
