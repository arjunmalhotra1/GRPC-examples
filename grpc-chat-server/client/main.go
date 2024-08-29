package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	chat "github.com/arjunmalhotra1/T-GRPC-2/grpc-chat-server/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Must have a url to connect to as the first argument, and a username as the second argument")
		return
	}

	ctx := context.Background()
	conn, err := grpc.NewClient(os.Args[1], grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	chatClient := chat.NewChatClient(conn)
	stream, err := chatClient.Chat(ctx)
	if err != nil {
		panic(err)
	}

	waitc := make(chan struct{})

	// This part is for receiving a message.
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			} else if err != nil {
				panic(err)
			}
			fmt.Println(msg.User + ": " + msg.Message)
		}
	}()

	fmt.Println("Connection established, type \"quit\" or use ctrl+c to exit")
	// fmt.Scan() doens't scan the whole lines of case but bufio.NewScanner() scans the whole line of code.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "quit" {
			err := stream.CloseSend()
			if err != nil {
				panic(err)
			}
			// this breaks out of the for loop
			break
		}

		err := stream.Send(&chat.ChatMessage{
			User:    os.Args[2],
			Message: msg,
		})

		if err != nil {
			panic(err)
		}
	}

	<-waitc
}
