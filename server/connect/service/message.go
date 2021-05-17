package service

import (
	"fmt"
)

//流通的消息单元 content为真正发送的内容
type Message struct {
	// timeline对应的数据库id
	ID           int64       `json:"id"`
	FromConnID   int64       `json:"from_conn_id"`
	ToConnID     int64       `json:"to_conn_id"`
	FromUID      string      `json:"from_uid"`
	ToID         string      `json:"to_id"`
	ReceiverType int         `json:"receiver_type"`
	Content      interface{} `json:"content"`
	Type         int         `json:"type"`
}

func NewMessage(client Client) *Message {
	return &Message{
		FromConnID: client.GetConnId(),
	}
}

func (m *Message) ToByte() []byte {
	return []byte(m.ToString())
}
func (m *Message) ToString() string {
	return fmt.Sprintf("%s", m.Content)
}
