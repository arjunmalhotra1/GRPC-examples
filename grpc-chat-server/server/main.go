package main

import (
	"fmt"
	"io"
	"net"
	"sync"

	chat "github.com/arjunmalhotra1/T-GRPC-2/grpc-chat-server/chat"
	"google.golang.org/grpc"
)

/*
Chat server receives a message from a client and then sends it to all the connections it is connected to at the time.
*/

type Connection struct {
	conn chat.Chat_ChatServer
	send chan *chat.ChatMessage
	quit chan struct{}
}

func NewConnection(conn chat.Chat_ChatServer) *Connection {
	c := &Connection{
		conn: conn,
		send: make(chan *chat.ChatMessage),
		quit: make(chan struct{}),
	}

	go c.start()
	return c
}

func (c *Connection) Send(msg *chat.ChatMessage) {
	defer func() {
		recover()
	}()
	c.send <- msg
}

func (c *Connection) start() {
	running := true
	for running {
		select {
		case msg := <-c.send:
			c.conn.Send(msg) // Ignoring the error, they just
		case <-c.quit:
			running = false
		}
	}
}

func (c *Connection) Close() error {
	close(c.quit)
	close(c.send)
	return nil
}

func (c *Connection) GetMessages(broadcast chan<- *chat.ChatMessage) error {
	for {
		message, err := c.conn.Recv()
		if err == io.EOF {
			c.Close()
			return nil
		} else if err != nil {
			c.Close()
			return err
		}

		go func(msg *chat.ChatMessage) {
			select {
			case broadcast <- msg:
			case <-c.quit:
			}
		}(message)
	}
}

type ChatServer struct {
	broadcast   chan *chat.ChatMessage
	quit        chan struct{}
	connections []*Connection
	connLock    sync.Mutex
	chat.UnimplementedChatServer
}

func NewChatServer() *ChatServer {
	srv := &ChatServer{
		broadcast: make(chan *chat.ChatMessage),
		quit:      make(chan struct{}),
	}

	go srv.start()
	return srv
}

func (c *ChatServer) Close() error {
	close(c.quit)
	return nil
}

func (c *ChatServer) start() {
	running := true

	for running {
		select {
		case msg := <-c.broadcast:
			c.connLock.Lock()
			for _, v := range c.connections {
				go v.Send(msg)
			}
			c.connLock.Unlock()
		case <-c.quit:
			running = false
		}
	}
}

// The ChatServer is what implements the rpc call in chat.proto
func (c *ChatServer) Chat(stream chat.Chat_ChatServer) error {
	conn := NewConnection(stream)

	c.connLock.Lock()
	c.connections = append(c.connections, conn)
	c.connLock.Unlock()

	err := conn.GetMessages(c.broadcast)

	c.connLock.Lock()
	// I don't know why are we deleting the connection here.
	for i, v := range c.connections {
		if v == conn {
			c.connections = append(c.connections[:i], c.connections[i+1:]...)
		}
	}

	c.connLock.Unlock()
	return err
}

func main() {
	lst, err := net.Listen("tcp", ":8887")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()

	srvc := NewChatServer()
	chat.RegisterChatServer(grpcServer, srvc)

	fmt.Println("Now serving at post 8887")
	err = grpcServer.Serve(lst)
	if err != nil {
		panic(err)
	}

}
