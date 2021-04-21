package db

import (
	"OpenFaas-User/model"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type DalUser struct {
	a *AffairDB
}

func NewDalUser() *DalUser {
	return &DalUser{}
}

func NewDalUserWithAffair(a *AffairDB) *DalUser {
	return &DalUser{a}
}

func (d *DalUser) GetDBHandler() *gorm.DB {
	if d.a != nil {
		return d.a.tx
	}
	return DBConn.FindersIMDB
}

func (d *DalUser) CreateUser(username, password string) (user *model.User, err error) {
	db := d.GetDBHandler()
	user = &model.User{
		UserID:   uuid.NewV4(),
		UserName: username,
		Password: password,
		Status:   model.Normal,
	}
	err = db.Create(user).Error
	return
}

func (d *DalUser) GetUserByUsername(username string) (user *model.User, err error) {
	db := d.GetDBHandler()
	user = new(model.User)
	err = db.Model(&model.User{}).Where("username = ?", username).First(user).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return
}

func (d *DalUser) GetUserByUserId(userId string) (user *model.User, err error) {
	db := d.GetDBHandler()
	user = new(model.User)
	err = db.Model(&model.User{}).Where("user_id = ?", userId).First(user).Error
	return
}
