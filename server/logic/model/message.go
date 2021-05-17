package model

import (
	"time"
)

// Message struct is a row record of the message table in the finders_imdb database
type Message struct {
	MessageID int64  `gorm:"primary_key;AUTO_INCREMENT;column:message_id;type:bigint;" json:"message_id"` // 消息ID
	DeviceID  int64  `gorm:"column:device_id;type:bigint;" json:"device_id"`                              // 设备唯一标示
	FromUID   string `gorm:"column:from_uid;type:varchar;size:64;" json:"from_uid"`                       // 消息发送人user_id
	ToID      string `gorm:"column:to_id;type:varchar;size:64;" json:"to_id"`                             // 消息接受者user_id
	//ToGid        int64      `gorm:"column:to_gid;type:bigint;" json:"to_gid"`                                    // 群组ID
	ReceiverType int        `gorm:"column:receiver_type;type:int;" json:"receiver_type"` // 接受者类型 1.个人 2.群组
	Type         int        `gorm:"column:type;type:int;" json:"type"`                   // 消息类型：1.文本 2.语音 3.视频 4.表情 6.图片
	Status       int        `gorm:"column:status;type:int;" json:"status"`               // 消息状态：1.未发送   2.已送达 (1的话会检查还有哪些用户没有收到)
	Content      string     `gorm:"column:content;type:text;size:65535;" json:"content"` // 消息内容，根据消息类型有所不同，json串
	CreatedAt    time.Time  `gorm:"column:created_at;type:datetime;" json:"created_at"`  // 创建时间
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:datetime;" json:"updated_at"`  // 更新时间
	DeletedAt    *time.Time `gorm:"column:deleted_at;type:datetime;" json:"deleted_at"`  // 删除时间

}

func (m *Message) TableName() string {
	return "message"
}
