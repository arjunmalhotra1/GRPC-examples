package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	blogs "github.com/arjunmalhotra1/T-GRPC-2/grpc-CRUD/blogs"
	"github.com/google/uuid"
)

type BlogItem struct {
	ID       uuid.UUID
	AuthorID string
	Content  string
	Title    string
}

type BlogServiceServer struct {
	blogs.UnimplementedBlogServiceServer
	blogStore map[uuid.UUID]BlogItem
}

func (b *BlogServiceServer) CreateBlog(ctx context.Context, req *blogs.CreateBlogReq) (*blogs.CreateBlogRes, error) {
	blog := req.GetBlog()
	data := BlogItem{
		ID:       uuid.New(),
		AuthorID: blog.AuthorId,
		Title:    blog.Title,
		Content:  blog.Content,
	}

	b.blogStore[data.ID] = data
	blog.Id = data.ID.String()
	return &blogs.CreateBlogRes{Blog: blog}, nil
}

func (b *BlogServiceServer) ReadBlog(ctx context.Context, req *blogs.ReadBlogReq) (*blogs.ReadBlogRes, error) {
	s, _ := uuid.Parse(req.Blog.Id)
	blogitem := b.blogStore[s]

	res := &blogs.ReadBlogRes{
		Blog: &blogs.Blog{
			Id:       req.Blog.Id,
			AuthorId: blogitem.AuthorID,
			Title:    blogitem.Title,
			Content:  blogitem.Content,
		},
	}

	return res, nil
}

func (b *BlogServiceServer) UpdateBlog(ctx context.Context, req *blogs.UpdateBlogReq) (*blogs.UpdateBlogRes, error) {
	return &blogs.UpdateBlogRes{}, nil
}
func (b *BlogServiceServer) DeleteBlog(ctx context.Context, req *blogs.DeleteBlogReq) (*blogs.DeleteBlogRes, error) {
	return &blogs.DeleteBlogRes{}, nil
}

func (b *BlogServiceServer) ListBlogs(req *blogs.ListBlogsRequest, stream blogs.BlogService_ListBlogsServer) error {
	for _, v := range b.blogStore {

		stream.Send(&blogs.ListBlogsResponse{
			Blog: &blogs.Blog{
				Id:       v.ID.String(),
				AuthorId: v.AuthorID,
				Title:    v.Title,
				Content:  v.Content,
			},
		})
	}

	return nil
}

func main() {

	blogStore := make(map[uuid.UUID]BlogItem)

	listener, err := net.Listen("tcp", ":8456")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	blogServc := &BlogServiceServer{
		blogStore: blogStore,
	}
	blogs.RegisterBlogServiceServer(grpcServer, blogServc)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to Serve")
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	grpcServer.Stop()
	listener.Close()
}
