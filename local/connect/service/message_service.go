package service

import (
	"OpenFaas-Connect/pkg/dto"
	"OpenFaas-Connect/st"
	"fmt"
)

type MessageService struct{}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (m *MessageService) NotifyMessage(msgModel *dto.NotifyMessageModel) (err error) {
	msg := &Message{
		ID:           msgModel.ID,
		FromConnID:   msgModel.FromConnID,
		ToConnID:     msgModel.ToConnID,
		FromUID:      msgModel.FromUID,
		ToID:         msgModel.ToID,
		Content:      msgModel.Content,
		ReceiverType: msgModel.ReceiverType,
		Type:         dto.MESSAGE_TEXT_TYPE,
	}
	client := HubSrv.GetClient(msg.ToID)
	st.Debug(msg.Content)
	if client == nil {
		return fmt.Errorf("user not online")
	}
	err = HubSrv.SendMsgToClient(client, msg)
	return
}
