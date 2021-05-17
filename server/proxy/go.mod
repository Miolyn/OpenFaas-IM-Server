module OpenFaas-Proxy

go 1.15

replace handler/function => ./function

require (
	github.com/gin-gonic/gin v1.7.1
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.25.0
)
