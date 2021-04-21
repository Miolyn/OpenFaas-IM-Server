package api

import (
	"OpenFaas-Connect/pkg/dto"
	"OpenFaas-Connect/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Notify(c *gin.Context) {
	var (
		form dto.NotifyMessageModel
		err  error
	)
	err = c.ShouldBind(&form)
	if err != nil {
		ErrorHandler(c, fmt.Sprintf("bind json error"))
		return
	}
	srvProxy := service.NewMessageService()
	err = srvProxy.NotifyMessage(&form)
	if err != nil {
		ErrorHandler(c, fmt.Sprintf("send notify error:%v\n", err))
		return
	}
	DataHandler(c, true)
}
