package api

import (
	"OpenFaas-User/model"
	"OpenFaas-User/pkg/cache"
	"OpenFaas-User/service"
	"OpenFaas-User/st"
	"OpenFaas-User/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RegisterOrLogin(c *gin.Context) {
	type FORM struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var (
		err    error
		form   FORM
		user   *model.User
		token  string
		userId string
	)
	err = c.ShouldBind(&form)
	if err != nil {
		ErrorHandler(c, fmt.Sprintf("json请求解析出错, %v", err.Error()))
		return
	}
	srv := service.NewUserService()
	cacheSrv := cache.NewUserCacheService()
	user, err = cacheSrv.GetUserByUsernameOrUserId(form.Username)
	if err == nil {
		userId = user.UserID.String()
		st.Debug("use user redis")
	}
	if userId == "" {

		if form.Username == "" || form.Password == "" {
			ErrorHandler(c, fmt.Sprintf("用户名或密码为空"))
			return
		}
		user, err = srv.GetUserByUsername(form.Username)
		if err != nil {
			ErrorHandler(c, fmt.Sprintf("mysql出错, %v", err.Error()))
			return
		}
		if user == nil {
			user, err = srv.Register(form.Username, form.Password)
			if err != nil {
				ErrorHandler(c, fmt.Sprintf("mysql出错, %v", err.Error()))
				return
			}
		}
		userId = user.UserID.String()
		err = cacheSrv.SetUserByUsernameAndUserId(user)
		if err != nil {
			ErrorHandler(c, fmt.Sprintf("redis出错, %v", err.Error()))
			return
		}
	}

	token, err = srv.GetUserToken(userId)
	if err != nil {
		ErrorHandler(c, fmt.Sprintf("生成token出错"))
		return
	}
	data := make(gin.H)
	data["token"] = token
	data["user_id"] = user.UserID.String()
	DataHandler(c, data)
}

func CheckToken(c *gin.Context) {
	type FORM struct {
		Username string `json:"username"`
		UserId   string `json:"user_id"`
	}
	var (
		err       error
		token     string
		user      *model.User
		jwtClaims *util.JWTClaims
		form      *FORM
	)
	token = c.GetHeader("token")
	if token == "" {
		ErrorHandler(c, fmt.Sprint("token 为空"))
		return
	}
	jwt := util.NewJWT()
	jwtClaims, err = jwt.ParseToken(token)
	if err != nil {
		ErrorHandler(c, fmt.Sprint("token 转换出错"))
		return
	}
	srv := service.NewUserService()
	cacheSrv := cache.NewUserCacheService()
	user, err = cacheSrv.GetUserByUsernameOrUserId(jwtClaims.UserID)
	if err == nil {
		form = &FORM{
			Username: user.UserName,
			UserId:   user.UserID.String(),
		}
		DataHandler(c, form)
		return
	}
	user, err = srv.GetUserByUserId(jwtClaims.UserID)
	if err != nil {
		ErrorHandler(c, fmt.Sprintf("获取用户信息失败, %v", err.Error()))
		return
	}
	err = cacheSrv.SetUserByUsernameAndUserId(user)
	if err != nil {
		ErrorHandler(c, fmt.Sprintf("redis出错, %v", err.Error()))
		return
	}
	form = &FORM{
		Username: user.UserName,
		UserId:   user.UserID.String(),
	}
	DataHandler(c, form)
}
