# 获取对应架构的of-watchdog镜像
FROM --platform=${TARGETPLATFORM:-linux/amd64} ghcr.io/openfaas/of-watchdog:0.8.4 as watchdog
# 获取对应架构的golang运行环境，使用的是基于Linux alpine 打好的 golang 镜像，可以减少镜像的大小
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.15-alpine3.13 as build
# 设置相关的构建时变量，方便跨平台构建
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
# 安装git
RUN apk --no-cache add git
# 为了镜像瘦身，编译纯静态的go程序，将环境变量CGO_ENABLED设置为0来关闭使用CGO编译器
ENV CGO_ENABLED=0
# 将of-watchdog镜像中的watchdog可执行文件复制到当前惊险中
COPY --from=watchdog /fwatchdog /usr/bin/fwatchdog
# 将其设置为可执行文件
RUN chmod +x /usr/bin/fwatchdog
# 创建handler文件夹作为工作路径
RUN mkdir -p /go/src/handler
# 设置工作路径
WORKDIR /go/src/handler
# 将宿主机的目录下代码复制到当前工作路径下准备编译使用
COPY . .

# Run a gofmt and exclude all vendored code.
# 使用gofmt一键格式化代码
RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./function/vendor/*"))" || { echo "Run \"gofmt -s -w\" on your Golang code"; exit 1; }
# 设置包的版本管理方式环境变量
ARG GO111MODULE="on"
# 设置获取包的国内镜像源
ARG GOPROXY="https://goproxy.cn"
# 首先go mod vender将依赖复制到vendor文件夹下方便查找使用。 然后根据平台和架构进行静态go程序编译
RUN go mod vendor \
    && CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build --ldflags "-s -w" -a -installsuffix cgo -o handler .
#RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go test handler/function/... -cover
# 引入新一阶段的镜像，减少镜像大小
FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3.13
# Add non root user and certs
# 添加非root用户和cert证书
RUN apk --no-cache add ca-certificates \
    && addgroup -S app && adduser -S -g app app
# Split instructions so that buildkit can run & cache 
# the previous command ahead of time.
# 创建对应用户的用户目录准备将可执行文件放入用户目录中
RUN mkdir -p /home/app \
    && chown app /home/app
# 设置WORKDIR工作目录
WORKDIR /home/app
# 将上一阶段准备好的可执行文件以及watchdog等复制到该工作目录下
COPY --from=build /go/src/handler/handler    .
COPY --from=build /usr/bin/fwatchdog         .
COPY --from=build /go/src/handler/function/  .
# 改变目录所属用户
RUN chown -R app /home/app
# 设置用户
USER app
# 配置of-watchdog相关的环境变量
# 处理函数
ENV fprocess="./handler"
# 处理模式
ENV mode="http"
# 镜像监控端口
ENV http_port 8084
# 请求转发的路径
ENV upstream_url="http://127.0.0.1:8084"
# 是否使用前缀日志
ENV prefix_logs="false"
# 运行启动of-watchdog进行函数监控
CMD ["./fwatchdog"]
