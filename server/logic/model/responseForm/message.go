package responseForm

import "time"

type Paginator struct {
	Page  int64       `json:"page"`
	Size  int64       `json:"size"`
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}

type SimpleMessage struct {
	MessageID    int64     `json:"message_id"`
	DeviceID     int64     `json:"device_id"`
	FromUID      string    `json:"from_uid"`
	ToID         string    `json:"to_id"`
	ReceiverType int       `json:"receiver_type"`
	Type         int       `json:"type"`
	Status       int       `json:"status"`
	Content      string    `json:"content"`
	Seq          int64     `json:"seq"`
	CreatedAt    time.Time `json:"created_at"`
}
