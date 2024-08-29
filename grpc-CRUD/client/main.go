package main

import (
	"context"
	"fmt"

	blogs "github.com/arjunmalhotra1/T-GRPC-2/grpc-CRUD/blogs"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	ctx := context.Background()

	conn, err := grpc.NewClient("localhost:8456", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	blogsClient := blogs.NewBlogServiceClient(conn)
	blogsClient.CreateBlog(ctx, &blogs.CreateBlogReq{
		Blog: &blogs.Blog{
			Id:       uuid.New().String(),
			AuthorId: "123",
			Title:    "test1",
			Content:  "test1",
		},
	})

	stream, _ := blogsClient.ListBlogs(ctx, &blogs.ListBlogsRequest{})

	data, _ := stream.Recv()

	fmt.Println("list of all blogs: ", data)
}
