package function

import (
	"OpenFaas-Connect/pkg/logger"
	"OpenFaas-Connect/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	handler "github.com/openfaas/templates-sdk/go-http"
	"net/http"
	"net/url"
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

func getQuery(queryString string, query string) string {
	values, _ := url.ParseQuery(queryString)
	return values[query][0]
}

func WSConnect(c *gin.Context) {
	var (
		err    error
		wsConn *websocket.Conn
	)

	userId, _ := c.GetQuery("user_id")
	deviceId, _ := strconv.ParseInt(userId, 10, 64)
	//token := c.GetHeader("token")
	token := c.Query("token")
	if token == "" {
		logger.Logger.Errorf("校验失败 token为空")
		return
		//wsConn.Close()
	}
	connId := deviceId
	wsConn, err = upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Logger.Errorf("获取链接失败 :%v", err)
		return
	}
	client := service.NewClient(service.HubSrv.Context(), wsConn)
	client.UserId = userId
	client.ConnID = connId
	client.DeviceId = deviceId
	client.Hub = service.HubSrv
	service.HubSrv.Register(client)
	// defer client.Close()
	go client.Start()
}

// Handle a function invocation
func Handle(resp http.ResponseWriter, req *http.Request) (handler.Response, error) {
	var (
		err    error
		wsConn *websocket.Conn
	)

	userId := getQuery(req.URL.RawQuery, "user_id")
	deviceId, _ := strconv.ParseInt(userId, 10, 64)
	wsConn, err = upgrader.Upgrade(resp, req, nil)
	if err != nil {
		return handler.Response{
			Body:       []byte("error"),
			StatusCode: http.StatusOK,
		}, err
	}
	client := service.NewClient(service.HubSrv.Context(), wsConn)
	client.UserId = userId
	client.ConnID = deviceId
	client.DeviceId = deviceId
	client.Hub = service.HubSrv
	fmt.Println(client.UserId, client.ConnID)

	service.HubSrv.Register(client)
	// defer client.Close()
	go client.Start()
	return handler.Response{
		Body:       []byte("ok"),
		StatusCode: http.StatusOK,
	}, err
}
