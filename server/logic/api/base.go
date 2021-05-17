package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseJson struct {
	Msg   string      `json:"msg"`
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Page  int64       `json:"page"`
	Size  int64       `json:"size"`
	Total int64       `json:"total"`
}

const (
	ResponseOK   = 0
	ResponseFail = -1
)

func DataHandler(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseJson{Data: data, Code: ResponseOK})
}

func ErrorHandler(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, &ResponseJson{Msg: msg, Data: false, Code: ResponseFail})
}
