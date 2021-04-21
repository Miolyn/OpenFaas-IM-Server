package api

import (
	"OpenFaaS-Logic/pkg/dto"
	"OpenFaaS-Logic/pkg/logger"
	"OpenFaaS-Logic/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

func SendMessage(c *gin.Context) {
	var (
		err error
		msg *dto.MessageModel
	)
	msg = new(dto.MessageModel)
	err = c.ShouldBind(msg)
	if err != nil {
		logger.Logger.Debug(err)
		ErrorHandler(c, fmt.Sprintf("json绑定失败"))
		return
	}
	if msg.Check(c) {
		return
	}
	if msg.FromUID == msg.ToID {
		ErrorHandler(c, fmt.Sprint("逻辑错误"))
		return
	}
	//srvProxy := service.MessageService{}
	srvProxy := service.NewMessageService()
	err = srvProxy.CreateMessage(msg)
	if err != nil {
		ErrorHandler(c, fmt.Sprintf("message 存储mysql失败"))
		return
	}
	err = srvProxy.SendMessage(msg)
	if err != nil {
		ErrorHandler(c, fmt.Sprintf("消息队列发送失败"))
		return
	}
	DataHandler(c, true)
}

//
//func GetUnReadMessage(c *gin.Context) {
//	var (
//		err    error
//		userID string
//		page   int64
//		form   responseForm.Paginator
//	)
//	userID = c.Query("user_id")
//	page = com.StrTo(c.DefaultQuery("page", "1")).MustInt64()
//	if userID == "" {
//		ErrorHandler(c, e.INFO_ERROR)
//		return
//	}
//	base := service.NewBase()
//	if base.AffairImInit() {
//		ErrorHandler(c, e.MYSQL_AFFAIR_ERROR)
//		return
//	}
//	defer base.AffairBegin(c)()
//
//	srvProxy := service.NewMessageService()
//	srvProxy.InitAffair(base)
//	fmt.Println(userID)
//	form, err = srvProxy.GetUnRead(userID, page, 8)
//	if err != nil {
//		base.AffairRollback()
//		ErrorHandler(c, fmt.Errorf("获取未读消息失败:%v", err).Error())
//		return
//	}
//	if base.AffairFinished() {
//		ErrorHandler(c, e.MYSQL_AFFAIR_ERROR)
//		return
//	}
//	PageHandler(c, form)
//}
//
//func GetTimeline(c *gin.Context) {
//	var (
//		err        error
//		page       int64
//		timelineId int64
//		seq        int64
//		form       responseForm.Paginator
//	)
//	page = com.StrTo(c.DefaultQuery("page", "1")).MustInt64()
//	timelineIdStr := c.Query("timeline_id")
//	seqStr := c.Query("seq")
//	if timelineIdStr == "" || seqStr == "" {
//		ErrorHandler(c, e.INFO_ERROR)
//		return
//	}
//	timelineId = com.StrTo(timelineIdStr).MustInt64()
//	seq = com.StrTo(seqStr).MustInt64()
//	base := service.NewBase()
//	if base.AffairImInit() {
//		ErrorHandler(c, e.MYSQL_AFFAIR_ERROR)
//		return
//	}
//
//	defer base.AffairBegin(c)()
//	srvProxy := service.NewMessageService()
//	srvProxy.InitAffair(base)
//	form, err = srvProxy.GetTimelineFormBySeq(timelineId, seq, page, 8)
//	if err != nil {
//		base.AffairRollback()
//		ErrorHandler(c, fmt.Errorf("获取timeline消息失败:%v", err).Error())
//		return
//	}
//	if base.AffairFinished() {
//		ErrorHandler(c, e.MYSQL_AFFAIR_ERROR)
//		return
//	}
//	PageHandler(c, form)
//}
