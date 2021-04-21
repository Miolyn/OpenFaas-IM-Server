package api

import (
	"OpenFaas-Connect/pkg/dto"
	"OpenFaas-Connect/pkg/logger"
	"OpenFaas-Connect/service"
	"OpenFaas-Connect/st"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	//定义读写缓冲区大小
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {

		return true
	},
}

func WSConnect(c *gin.Context) {
	var (
		err    error
		wsConn *websocket.Conn
		resp   *dto.BindResponseJson
	)

	token := c.Query("token")
	srv := service.NewRegisterService()
	resp, err = srv.Register(token)
	if err != nil {
		logger.Logger.Errorf("获取链接失败 :%v", err)
		return
	}
	userId := resp.Data.UserId

	wsConn, err = upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Logger.Errorf("获取链接失败 :%v", err)
		return
	}
	client := service.NewClient(service.HubSrv.Context(), wsConn)
	client.UserId = userId
	client.ConnID, _ = strconv.ParseInt("123", 10, 64)
	client.DeviceId = client.ConnID
	client.Username = resp.Data.Username
	client.Hub = service.HubSrv
	service.HubSrv.Register(client)
	// defer client.Close()
	st.Debug(client.UserId, " connect")
	go client.Start()
}
