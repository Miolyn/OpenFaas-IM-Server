package db

import (
	"OpenFaaS-Logic/model"
	"github.com/jinzhu/gorm"
)

type DalMessage struct {
	a *AffairDB
}

func NewDalMessage() *DalMessage {

	return &DalMessage{}
}
func NewDalMessageAffair(a *AffairDB) *DalMessage {

	return &DalMessage{a}
}
func (d *DalMessage) GetDBHandler() *gorm.DB {
	if d.a != nil {
		return d.a.tx
	}
	return DBConn.FindersIMDB
}

func (d *DalMessage) CreateMessage(message *model.Message) (err error) {
	//message = new(model.Message)
	db := d.GetDBHandler().Debug().Table(message.TableName())
	err = db.Create(message).Error
	return err
}

func (d *DalMessage) GetMessageByMessageID(messageID int64) (message *model.Message, err error) {
	db := d.GetDBHandler()
	message = new(model.Message)
	db = db.Table(message.TableName())
	err = db.Where("message_id = ?", messageID).First(message).Error
	return
}
