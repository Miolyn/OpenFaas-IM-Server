package util

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func ToIntP(num int) *int {
	np := &num
	return np
}

func ToInt64P(num int64) *int64 {
	np := &num
	return np
}

func ToStrP(str string) *string {
	s := &str
	return s
}

func QueryInt(c *gin.Context, name string) (res int64) {
	tmp := c.Query(name)
	res, _ = strconv.ParseInt(tmp, 10, 64)
	return
}
