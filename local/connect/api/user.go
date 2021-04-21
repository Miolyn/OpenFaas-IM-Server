package api

import (
	"OpenFaas-Connect/pkg/dto"
	"OpenFaas-Connect/service"
	"github.com/gin-gonic/gin"
)

func GetUsersOnline(c *gin.Context) {
	var (
		form []dto.UserForm
	)
	form = service.HubSrv.GetAllOnline()
	DataHandler(c, form)
}
