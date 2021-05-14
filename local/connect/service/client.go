package service

import (
	"OpenFaas-Connect/pkg/logger"
	"context"
	"github.com/gorilla/websocket"
	"net"
	"sync"
	"time"
)

const (
	//Time allowed to write message to the peer
	writeWait = 10 * time.Second
	//writeWait = 8 * time.Second

	//Time allowed to read the next pong message from the peer
	//pongWait = 10 * time.Second
	pongWait = 50 * time.Second

	//Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	//Maximum message size allowed from peer
	maxMessageSize int64 = 1024
)

type Client interface {
	Close()
	GetAddr() string
	GetUserId() string
	GetUsername() string
	GetConnId() int64
	GetDeviceId() int64
	GetToken() string
	SendMsg(*Message)
	GetConn() *websocket.Conn
	Read() <-chan *Message
	Start() //启动读写 goroutine
	// Read
}

type client struct {
	Addr          string //客户端地址
	UserId        string //用户ID uuid
	Username      string
	ConnID        int64  //连接Id 全局唯一
	DeviceId      int64  //设备Id
	Token         string //登陆令牌
	SendChan      chan *Message
	ReadChan      chan *Message
	Conn          *websocket.Conn //用户连接
	LoginTime     int64           //unix时间戳
	HeartBeatTime int64           //心跳时间
	// sync.Once能确保实例化对象Do方法在多线程环境只运行一次,内部通过互斥锁实现
	CloseOnce  sync.Once
	Ctx        context.Context
	CancelFunc context.CancelFunc

	Hub Hub //属于的处理器
}

func NewClient(ctx context.Context, conn *websocket.Conn) *client {
	ctx, cancelFunc := context.WithCancel(ctx)
	return &client{
		LoginTime:  time.Now().Unix(),
		SendChan:   make(chan *Message, 1024),
		ReadChan:   make(chan *Message, 1024),
		Conn:       conn,
		Addr:       conn.RemoteAddr().String(),
		Ctx:        ctx,
		CancelFunc: cancelFunc,
	}
}

func (c *client) Start() {
	go c.readPump(c.Ctx)
	go c.writePump(c.Ctx)
	for {
		select {
		case message, isOpen := <-c.Read():
			if isOpen && message != nil {
				logger.Logger.Infof("用户<%v> 收到消息：%s", message.ToConnID, message.Content)
			}
		// 关闭client
		case <-c.Ctx.Done():
			logger.Logger.Errorf("Client %v 正常关闭 err:%v", c.ConnID, c.Ctx.Err())
			c.Close()
			return
		}
	}
}

func (c *client) Close() {
	c.CloseOnce.Do(func() {
		logger.Logger.Noticef("注销用户 %v \n", c.ConnID)
		close(c.SendChan)
		close(c.ReadChan)
		c.CancelFunc()
		c.Hub.UnRegister(c)
		c.Conn.Close()
	})
	//defer c.Conn.Close()

}

func (c *client) GetAddr() string {
	return c.Addr
}
func (c *client) GetUserId() string {
	return c.UserId
}
func (c *client) GetUsername() string {
	return c.Username
}
func (c *client) GetConnId() int64 {
	return c.ConnID
}
func (c *client) GetDeviceId() int64 {
	return c.DeviceId
}
func (c *client) GetToken() string {
	return c.Token
}
func (c *client) SendMsg(msg *Message) {
	c.SendChan <- msg
}
func (c *client) GetConn() *websocket.Conn {
	return c.Conn
}
func (c *client) Read() <-chan *Message {
	return c.ReadChan
}



// readPump pumps messages from the websocket connection to the HubSrv
func (c *client) readPump(ctx context.Context) {
	defer c.Close()
	// 设置pong的deadline 若超时没有接收到pong则自动断开连接
	c.GetConn().SetReadDeadline(time.Now().Add(pongWait))
	c.GetConn().SetReadLimit(maxMessageSize)
	c.GetConn().SetPongHandler(func(string) error {
		//logger.Logger.Infof("用户<%v> Pong hanlder...", c.ConnID)
		c.HeartBeatTime = time.Now().Unix()
		c.GetConn().SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, msg, err := c.Conn.ReadMessage()
			logger.Logger.Infof("connId: %s, receive msg:%s", c.ConnID, msg)
			if err != nil {
				//if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
					logger.Logger.Errorf("error: %v", err)
					return
				}
				//if err, ok := err.(*net.OpError); ok && err.{
				if err, ok := err.(net.Error); ok && err.Timeout() {
					logger.Logger.Info("read timeout")
					return
				}
				logger.Logger.Infof("receive  error %v", err)
				break
			}
			message := NewMessage(c)
			message.Content = msg
			message.ToConnID = c.GetConnId()
			c.ReadChan <- message
		}

	}

}

// writePump pumps messages from the HubSrv to the websocket connection
func (c *client) writePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case message, isOpen := <-c.SendChan:
			if isOpen {
				//c.GetConn().SetWriteDeadline(time.Now().Add(writeWait))
				c.GetConn().SetWriteDeadline(time.Now().Add(pongWait))
				//if err := c.GetConn().WriteJSON(map[string]interface{}{"data": message.ToString()}); err != nil {
				if err := c.GetConn().WriteJSON(map[string]interface{}{"data": message}); err != nil {
					logger.Logger.Errorf("发送信息失败 %v", err)
				}
			}
		case <-ticker.C:
			//logger.Logger.Info("向用户发送ping")
			//c.GetConn().SetWriteDeadline(time.Now().Add(writeWait))
			c.GetConn().SetWriteDeadline(time.Now().Add(pongWait))
			if err := c.GetConn().WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Logger.Errorf("给用户<%v>发送ping 失败 ：%v", c.GetConnId(), err)
				return
			}
		case <-ctx.Done():
			return
		}
	}

}
