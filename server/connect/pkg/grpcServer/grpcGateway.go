package grpcServer

import (
	"OpenFaas-Connect/pkg/grpcServer/pb"

	"context"
	"flag"
	"fmt"
	//"github.com/rs/cors"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net/http"
)

func run(port string, echoEndpoint *string) error {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	//handler := cors.AllowAll().Handler(mux)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterMessageHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
	//cors.Wrap(mux, ExampleCORSOptions()...)

	if err != nil {
		return err
	}
	//handler := cors.AllowAll().Handler(mux)
	handler := mux
	return http.ListenAndServe(":"+port, handler)
}

func GRPCGateway(port, grpcServerPort string) {
	var (
		echoEndpoint = flag.String("echo_endpoint", "localhost:"+grpcServerPort, "endpoint of YourService")
	)
	go func() {
		if err := run(port, echoEndpoint); err != nil {
			fmt.Print(err.Error())
		}
	}()

}
