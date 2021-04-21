package service

import (
	"OpenFaas-User/dal/db"
	"OpenFaas-User/model"
	"OpenFaas-User/util"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) GetUserByUsername(username string) (user *model.User, err error) {
	dal := db.NewDalUser()
	user, err = dal.GetUserByUsername(username)
	return
}

func (u *UserService) GetUserByUserId(userId string) (user *model.User, err error) {
	dal := db.NewDalUser()
	return dal.GetUserByUserId(userId)
}

func (u *UserService) Register(username, password string) (user *model.User, err error) {
	dal := db.NewDalUser()
	return dal.CreateUser(username, password)
}

func (u *UserService) GetUserToken(userId string) (token string, err error) {
	jwt := util.NewJWT()
	jwtClaims := util.JWTClaims{
		MapClaims: nil,
		UserID:    userId,
	}
	token, err = jwt.GenerateToken(jwtClaims)
	return
}
