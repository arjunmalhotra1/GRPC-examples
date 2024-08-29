package main

import (
	"context"
	"fmt"

	echo "github.com/arjunmalhotra1/T-GRPC-2/01-proto/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	conn, err := grpc.NewClient("localhost:8887", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	echoClient := echo.NewEchoServerClient(conn)
	resp, err := echoClient.Echo(ctx, &echo.EchoRequest{
		Message: "ClientMessage: Hello World!",
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Got from Server: ", resp.Response)
}
