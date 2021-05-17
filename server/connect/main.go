package main

import (
	"OpenFaas-Connect/api"
	function "OpenFaas-Connect/function"
	"OpenFaas-Connect/pkg/grpcServer"
	"OpenFaas-Connect/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	handler "github.com/openfaas/templates-sdk/go-http"
)

var (
	acceptingConnections int32
)

const defaultTimeout = 10 * time.Second
const port = 8082
const grpcPort = "8085"
const gatewayPort = "8086"

// 121.36.35.85
func main() {
	readTimeout := parseIntOrDurationValue(os.Getenv("read_timeout"), defaultTimeout)
	writeTimeout := parseIntOrDurationValue(os.Getenv("write_timeout"), defaultTimeout)
	service.HubSrv = service.NewHub(context.Background())
	go service.HubSrv.Run()
	g := gin.Default()
	g.GET("/ws", api.WSConnect)
	g.GET("/online", api.GetUsersOnline)
	g.POST("/notify", api.Notify)
	grpcServer.StartGRPCServer(grpcPort)
	grpcServer.GRPCGateway(gatewayPort, grpcPort)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20, // Max header of 1MB
		Handler:        g,
	}
	fmt.Println("listen on :", port)
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

func makeRequestHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input []byte

		if r.Body != nil {
			defer r.Body.Close()

			bodyBytes, bodyErr := ioutil.ReadAll(r.Body)

			if bodyErr != nil {
				log.Printf("Error reading body from request.")
			}

			input = bodyBytes
		}

		req := handler.Request{
			Body:        input,
			Header:      r.Header,
			Method:      r.Method,
			QueryString: r.URL.RawQuery,
		}
		req.WithContext(r.Context())

		result, resultErr := function.Handle(w, r)

		if result.Header != nil {
			for k, v := range result.Header {
				w.Header()[k] = v
			}
		}

		if resultErr != nil {
			log.Print(resultErr)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			if result.StatusCode == 0 {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(result.StatusCode)
			}
		}

		w.Write(result.Body)
	}
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
