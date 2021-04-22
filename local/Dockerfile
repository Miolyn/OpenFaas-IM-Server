FROM golang:latest

MAINTAINER Miolyn "tju_liyining@foxmail.com"

WORKDIR $GOPATH/src/github.com/OpenFaas-Total

COPY . $GOPATH/src/github.com/OpenFaas-Total

ARG GO111MODULE="on"
ARG GOPROXY="https://goproxy.cn"
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
WORKDIR $GOPATH/src/github.com/OpenFaas-Total/connect

RUN go mod vendor && GOPROXY="https://goproxy.cn" GO111MODULE=on  go build -o main .

WORKDIR $GOPATH/src/github.com/OpenFaas-Total/logic

RUN go mod vendor && GOPROXY="https://goproxy.cn" GO111MODULE=on  go build -o main .

WORKDIR $GOPATH/src/github.com/OpenFaas-Total/proxy

RUN go mod vendor && GOPROXY="https://goproxy.cn" GO111MODULE=on  go build -o main .

WORKDIR $GOPATH/src/github.com/OpenFaas-Total/user

RUN go mod vendor && GOPROXY="https://goproxy.cn" GO111MODULE=on  go build -o main .


EXPOSE 8082 8083 8084 8085 8086 8087

WORKDIR $GOPATH/src/github.com/OpenFaas-Total

#ADD run.sh .
#CMD ["/bin/sh","./run.sh"]
#RUN nohup ./connect/main > /tmp/connect.log 2>&1 &
#RUN nohup ./logic/main > /tmp/logic.log 2>&1 &
#RUN nohup ./user/main > /tmp/user.log 2>&1 &
#RUN nohup ./proxy/main > /tmp/proxy.log 2>&1 &
#--network =“host”
ENTRYPOINT  ["./run.sh"]