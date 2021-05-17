package grpcServer

import (
	"OpenFaas-Connect/pkg/dto"
	"OpenFaas-Connect/pkg/grpcServer/pb"
	"OpenFaas-Connect/service"
	"OpenFaas-Connect/st"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedMessageServer
}

func (s *server) Notify(ctx context.Context, request *pb.NotifyRequest) (*pb.Response, error) {
	//return nil, status.Errorf(codes.Unimplemented, "method Notify not implemented")
	var (
		form    dto.NotifyMessageModel
		err     error
		content dto.NotifyContent
	)
	st.Debug(*request)
	st.Debug(request.Content)
	json.Unmarshal([]byte(request.Content), &content)
	st.Debug(content)
	form = dto.NotifyMessageModel{
		ID:           request.Id,
		OpCode:       request.OpCode,
		ToConnID:     request.ToConnID,
		FromConnID:   request.FromConnId,
		FromUID:      request.FromUid,
		ToID:         request.ToId,
		ReceiverType: int(request.ReceiverType),
		Content:      content,
	}
	srvProxy := service.NewMessageService()
	err = srvProxy.NotifyMessage(&form)
	if err != nil {
		return &pb.Response{
			Msg:   fmt.Sprintf("send notify error:%v\n", err),
			Code:  -1,
			Data:  "",
			Page:  0,
			Size:  0,
			Total: 0,
		}, err
	}
	return &pb.Response{
		Msg:   "ok",
		Code:  0,
		Data:  "",
		Page:  0,
		Size:  0,
		Total: 0,
	}, nil
}
func (s *server) GetUserOnline(context.Context, *empty.Empty) (*pb.UserOnlineResponseForm, error) {
	var (
		form []dto.UserForm
		data []*pb.UserForm
	)
	form = service.HubSrv.GetAllOnline()
	for _, client := range form {
		data = append(data, &pb.UserForm{
			UserId:   client.UserId,
			Username: client.Username,
		})
	}
	return &pb.UserOnlineResponseForm{
		Msg:   "ok",
		Code:  0,
		Data:  data,
		Page:  0,
		Size:  0,
		Total: 0,
	}, nil
}

func StartGRPCServer(port string) {
	// 要监听的协议和端口
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 实例化gRPC server结构体
	s := grpc.NewServer()

	// 服务注册
	pb.RegisterMessageServer(s, &server{})

	log.Println("开始监听，等待远程调用...")

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

}
