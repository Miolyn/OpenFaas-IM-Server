package model

import (
	"time"
)

// MessageAck struct is a row record of the message_ack table in the finders_imdb database
type MessageAck struct {
	ID        int64     `gorm:"primary_key;AUTO_INCREMENT;column:id;type:bigint;" json:"id"` // 数据库自增主键
	MessageID int64     `gorm:"column:message_id;type:bigint;" json:"message_id"`            // 消息ID
	UserID    string    `gorm:"column:user_id;type:varchar;size:64;" json:"user_id"`         // 接收到此消息的user_id
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;" json:"updated_at"`
}

func (m *MessageAck) TableName() string {
	return "message_ack"
}
