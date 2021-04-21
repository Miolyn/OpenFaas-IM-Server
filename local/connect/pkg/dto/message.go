package dto

import (
	"OpenFaas-Connect/pkg/e"
	"OpenFaas-Connect/pkg/response"
	"OpenFaas-Connect/st"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"time"
)

const (
	RECV_PERSON = 1 //个人
	RECV_GROUP  = 2 //群组

	TIMELINE_PERSON = 1 //个人
	TIMELINE_DOUBLE = 2 //单聊
	TIMELINE_GROUP  = 3 //群组
)

// message type
const (
	MESSAGE_TEXT_TYPE    = 1
	MESSAGE_VOICE_TYPE   = 2
	MESSAGE_VIDEO_TYPE   = 3
	MESSAGE_EMOJI_TYPE   = 4
	MESSAGE_PICTURE_TYPE = 6
)

// message status
const (
	MESSAGE_NOT_SEND = iota + 1
	MESSAGE_SEND
)

type MessageModel struct {
	MessageID int64  `json:"message_id"`                          // 消息ID
	DeviceID  int64  `json:"device_id" validate:"required,gte=0"` // 设备唯一标示
	FromUID   string `json:"from_uid" validate:"required,gte=0"`  // 消息发送人user_id
	ToID      string `json:"to_id" validate:"required,gte=0"`     // 消息接受者user_id
	//ToGid        int64     `json:"to_gid"`        // 群组ID
	ReceiverType int       `json:"receiver_type" validate:"required,gte=0"` // 接受者类型 1.个人 2.群组
	MessageType  int       `json:"message_type" validate:"required,gte=0"`  // 消息类型：1.文本 2.语音 3.视频 4.表情 6.图片
	Status       int       `json:"status"`                                  // 消息状态：1.未发送   2.已送达 (1的话会检查还有哪些用户没有收到)
	Content      string    `json:"content" validate:"required,gte=0"`       // 消息内容，根据消息类型有所不同，json串
	CreatedAt    time.Time `json:"created_at"`                              // 创建时间
}

func (u *MessageModel) Check(c *gin.Context) bool {
	validate := validator.New()
	err := validate.Struct(*u)
	if err != nil {
		st.DebugWithFuncName(err)
		response.ErrorHandler(c, e.INFO_ERROR)
		return true
	}
	return false
}

type TimeLineModel struct {
	*MessageModel `json:"message"`
	FromConnID    int64     `json:"from_conn_id"`
	ToConnID      int64     `json:"to_conn_id"` //发送给哪个connection
	ID            int64     `json:"id"`
	TimelineID    int64     `json:"timeline_id"` // timeline ID 一个群，二个人，个人都分别有一个timeline_id
	ObjectIds     string    `json:"object_ids"`
	Seq           int64     `json:"seq"`           // 消息序列，一个timeline下的seq从0递增。
	TimeLineType  int       `json:"timeline_type"` // timeline类型，1.个人 2.单聊 3.群聊
	Status        int       `json:"status"`        // timeline状态，1.未读 2.已读
	CreatedAt     time.Time `json:"created_at"`
}

// op code
const (
	OPCODE_NEW_MSG = 1 + iota
)

type NotifyMessageModel struct {
	ID           int64       `json:"id"`      // 用于通知消息的对应，发送给前端这个id，前端响应已收到的时候也返回这个id
	OpCode       int64       `json:"op_code"` // 表示是哪种类型的通知
	ToConnID     int64       `json:"to_conn_id"`
	FromConnID   int64       `json:"from_conn_id"`
	FromUID      string      `json:"from_uid"`
	ToID         string      `json:"to_id"`
	ReceiverType int         `json:"receiver_type"` // 接受者类型 1.个人 2.群组
	Content      interface{} `json:"content"`
}

type NotifyContent struct {
	Content     string    `json:"content"` // new message notify
	MessageType int       `json:"message_type"`
	TimelineID  int64     `json:"timeline_id"`
	Seq         int64     `json:"seq"`
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
}
