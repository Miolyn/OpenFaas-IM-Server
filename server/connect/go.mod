module OpenFaas-Connect

go 1.15

replace handler/function => ./function

require (
	github.com/dan-compton/grpc-gateway-cors v0.0.0-20180203205722-9ab0d0703923
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/garyburd/redigo v1.6.2
	github.com/gin-gonic/gin v1.7.1
	github.com/go-playground/validator/v10 v10.5.0
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/imroc/req v0.3.0
	github.com/kr/text v0.2.0 // indirect
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/openfaas/templates-sdk v0.0.0-20200723110415-a699ec277c12
	github.com/rs/cors v1.7.0
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	google.golang.org/genproto v0.0.0-20210416161957-9910b6c460de
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
