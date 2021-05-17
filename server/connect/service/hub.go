package service

import (
	"OpenFaas-Connect/pkg/cache"
	"OpenFaas-Connect/pkg/dto"
	"OpenFaas-Connect/pkg/logger"
	"context"
	"fmt"
	"sync"
)

var (
	HubSrv Hub
)

type Hub interface {
	Run()              //启动
	Close()            //关闭
	AllClient() int64  //连接的clients num
	Register(Client)   //注册client
	UnRegister(Client) //注销client
	GetClient(userId string) Client
	SendMsgToClient(Client, *Message) error
	CheckAlive(Client) error
	Errors() <-chan error
	BroadCast(*Message)       //给连接的client 广播
	Context() context.Context //获取上下文
	ChildContext() (context.Context, context.CancelFunc)
	GetAllOnline() []dto.UserForm
}

//总处理器 实现	Hub 接口
type hubCli struct {
	//客户端列表
	clients map[string]Client //user_id - client
	ctx     context.Context
	//注册channel
	register chan Client
	//注销channel
	unregister chan Client

	//广播
	broadcast chan *Message
	errors    chan error
	close     chan error

	cLock sync.RWMutex //clients_map读写锁

	allClient int64 //总连接
}

func NewHub(ctx context.Context) Hub {
	h := &hubCli{
		ctx:        ctx,
		clients:    make(map[string]Client),
		register:   make(chan Client, 1000),
		unregister: make(chan Client, 1000),
		broadcast:  make(chan *Message, 1024),
		errors:     make(chan error),
		close:      make(chan error),
	}
	return h
}

func (h *hubCli) Run() {
	fmt.Println("hub running")
	for {
		select {
		//用户注册
		case client := <-h.register:
			h.registerClient(client)
			cache.SetUser(client.GetUserId(), client.GetConnId())
		//用户注销
		case client := <-h.unregister:
			h.unRegisterClient(client)
			cache.DelUser(client.GetUserId())
		//广播消息
		case message := <-h.broadcast:
			for _, client := range h.clients {
				if err := h.CheckAlive(client); err == nil {
					h.SendMsgToClient(client, message)
				}
			}
		//退出
		case done := <-h.close:
			if done != nil {
				logger.Logger.Infof("Hub 异常退出 %v ", done)
			} else {
				logger.Logger.Infof("Hub 正常退出")
			}
			close(h.close)
			return
		//一些程序异常 not fatal
		case err := <-h.Errors():
			logger.Logger.Errorf("Hub Error %v", err)
		case <-h.ctx.Done():
			logger.Logger.Errorf("context cancel")
			return
		}

	}
}

func (h *hubCli) Close() {
	close(h.broadcast)
	close(h.errors)
	close(h.register)
	close(h.unregister)
	h.close <- nil
}

func (h *hubCli) AllClient() int64 {
	return h.allClient
}

func (h *hubCli) Register(client Client) {
	h.register <- client
}
func (h *hubCli) UnRegister(client Client) {
	h.unregister <- client
}

func (h *hubCli) GetClient(userId string) Client {
	h.cLock.RLock()
	defer h.cLock.RUnlock()
	if client, ok := h.clients[userId]; ok {
		return client
	} else {
		return nil
	}
}

func (h *hubCli) SendMsgToClient(client Client, msg *Message) error {
	client.SendMsg(msg)
	return nil
}

func (h *hubCli) CheckAlive(client Client) error {
	if _, ok := h.clients[client.GetUserId()]; !ok {
		return fmt.Errorf("用户不在线")
	}
	return nil
}

func (h *hubCli) Errors() <-chan error {
	return h.errors
}
func (h *hubCli) BroadCast(msg *Message) {
	h.broadcast <- msg
}
func (h *hubCli) Context() context.Context {
	return h.ctx
}

func (h *hubCli) ChildContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(h.ctx)
}

//注册用户
func (h *hubCli) registerClient(client Client) {
	h.cLock.Lock()
	defer h.cLock.Unlock()

	h.clients[client.GetUserId()] = client
	h.allClient++
	logger.Logger.Infof("用户注册：%v  当前用户数量：%v", client.GetAddr(), len(h.clients))
}

//注销
func (h *hubCli) unRegisterClient(client Client) {
	h.cLock.RLock()
	defer h.cLock.RUnlock()
	if _, ok := h.clients[client.GetUserId()]; ok {
		delete(h.clients, client.GetUserId())
		h.allClient--
		logger.Logger.Infof("用户<%v>注销：%v  当前用户数量：%v", client.GetUserId(), client.GetAddr(), len(h.clients))
	}
}

func (h *hubCli) GetAllOnline() (form []dto.UserForm) {
	for _, client := range h.clients {
		form = append(form, *&dto.UserForm{
			Username: client.GetUsername(),
			UserId:   client.GetUserId(),
		})
	}
	return
}
