package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseJson struct {
	Msg   string      `json:"msg"`
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Page  int64       `json:"page"`
	Size  int64       `json:"size"`
	Total int64       `json:"total"`
}

func DataHandler(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseJson{Data: data, Code: 0})
}

func ErrorHandler(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, &ResponseJson{Msg: msg, Data: false, Code: -1})
}
