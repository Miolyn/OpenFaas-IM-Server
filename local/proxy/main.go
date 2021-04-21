package main

import (
	"OpenFaas-Proxy/middleware"
	"OpenFaas-Proxy/pkg/dto"
	"OpenFaas-Proxy/pkg/pb"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	acceptingConnections int32
)

//
const defaultTimeout = 10 * time.Second
const port = 8084

const userServerHost = "127.0.0.1:8083"

const logicServerHost = "127.0.0.1:8081"

const connectServerHost = "127.0.0.1:8086"

const connectGRPCServerHost = "127.0.0.1:8085"
const baseUrl = ""
const userBase = ""
const proxyBase = ""
const logicBase = ""

var client pb.MessageClient

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

func initGRPC() {
	conn, err := grpc.Dial(connectGRPCServerHost, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	client = pb.NewMessageClient(conn)
}

// 119.3.229.43
func main() {
	readTimeout := parseIntOrDurationValue(os.Getenv("read_timeout"), defaultTimeout)
	writeTimeout := parseIntOrDurationValue(os.Getenv("write_timeout"), defaultTimeout)
	initGRPC()
	g := gin.New()
	// logic
	g.POST("send_msg", func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = logicServerHost
			req.URL.Path = baseUrl + logicBase + "/send_msg"
			req.Host = logicServerHost
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	})
	g.Use(middleware.Cors())

	// connect
	g.GET("/online", func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = connectServerHost
			req.URL.Path = "/online"
			req.Host = connectServerHost
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	})
	g.POST("/notify", func(c *gin.Context) {
		var form dto.NotifyMessageModel
		err := c.ShouldBind(&form)
		if err != nil {
			ErrorHandler(c, "bing error")
			return
		}
		data, _ := json.Marshal(form.Content)
		fmt.Println(string(data))
		_, err = client.Notify(context.Background(), &pb.NotifyRequest{
			Id:           form.ID,
			OpCode:       form.OpCode,
			ToConnID:     form.ToConnID,
			FromConnId:   form.FromConnID,
			FromUid:      form.FromUID,
			ToId:         form.ToID,
			ReceiverType: int32(form.ReceiverType),
			Content:      string(data),
		})
		if err != nil {
			ErrorHandler(c, "call grpc error")
			return
		}

		DataHandler(c, "")
	})
	// user
	g.POST("login", func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = userServerHost
			req.URL.Path = baseUrl + userBase + "/login"
			req.Host = userServerHost

		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	})
	g.GET("check_token", func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = userServerHost
			req.URL.Path = baseUrl + userBase + "/check_token"
			req.Host = userServerHost
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20, // Max header of 1MB
		Handler:        g,
	}
	fmt.Println("listen on ", port)
	//http.HandleFunc("/ws", makeRequestHandler())
	listenUntilShutdown(s, writeTimeout)
}

func listenUntilShutdown(s *http.Server, shutdownTimeout time.Duration) {
	idleConnsClosed := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM)

		<-sig

		log.Printf("[entrypoint] SIGTERM received.. shutting down server in %s\n", shutdownTimeout.String())

		<-time.Tick(shutdownTimeout)

		if err := s.Shutdown(context.Background()); err != nil {
			log.Printf("[entrypoint] Error in Shutdown: %v", err)
		}

		log.Printf("[entrypoint] No new connections allowed. Exiting in: %s\n", shutdownTimeout.String())

		<-time.Tick(shutdownTimeout)

		close(idleConnsClosed)
	}()

	// Run the HTTP server in a separate go-routine.
	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("[entrypoint] Error ListenAndServe: %v", err)
			close(idleConnsClosed)
		}
	}()

	atomic.StoreInt32(&acceptingConnections, 1)

	<-idleConnsClosed
}

func parseIntOrDurationValue(val string, fallback time.Duration) time.Duration {
	if len(val) > 0 {
		parsedVal, parseErr := strconv.Atoi(val)
		if parseErr == nil && parsedVal >= 0 {
			return time.Duration(parsedVal) * time.Second
		}
	}

	duration, durationErr := time.ParseDuration(val)
	if durationErr != nil {
		return fallback
	}
	return duration
}
