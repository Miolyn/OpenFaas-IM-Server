package service

import (
	"OpenFaaS-Logic/dal/db"
	"OpenFaaS-Logic/pkg/dto"
	"testing"
)

func TestMessage(t *testing.T) {
	db.InitDB()
	// srvProxy := NEe()
	// dalProxy := db.NewDalTimelineMap()
	message := new(dto.MessageModel)
	message.FromUID = "elyar"
	message.ToID = "ablimit"
	message.ReceiverType = dto.RECV_PERSON

	res, err := MessageToTimeline(message)
	// m, err := dalProxy.GetTimeLineMap(2, "h14124k", "hellowrold")
	// userId := "1711306"
	// conn, err := srvProxy.GetConnID(nil, &userId)
	src := NewMessageService()
	// res[0].Content = "hello world"
	// fmt.Println("--------", res)
	err = src.SendMessage(message)
	t.Log(len(res), err)
}
