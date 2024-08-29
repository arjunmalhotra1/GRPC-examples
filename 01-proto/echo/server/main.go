package main

import (
	"context"
	"fmt"
	"net"

	echo "github.com/arjunmalhotra1/T-GRPC-2/01-proto/echo"
	"google.golang.org/grpc"
)

type EchoServer struct {
	echo.UnimplementedEchoServerServer
}

func (e *EchoServer) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{
		Response: "Anything: " + req.Message,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8887")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	server := &EchoServer{}

	echo.RegisterEchoServerServer(grpcServer, server)

	fmt.Println("Now serving at port 8887")
	err = grpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}
}
