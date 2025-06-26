package client

import (
	pb "blog-app/proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BlogClient struct {
	client pb.PostServiceClient
}

func NewBlogClient(addr string) (*BlogClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &BlogClient{
		client: pb.NewPostServiceClient(conn),
	}, nil
}

func (c *BlogClient) CreatePost(req *pb.CreatePostRequest) (*pb.Post, error) {
	ctx := context.Background()
	res, err := c.client.CreatePost(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *BlogClient) GetPost(req *pb.ReadPostRequest) (*pb.Post, error) {
	ctx := context.Background()
	res, err := c.client.ReadPost(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *BlogClient) UpdatePost(req *pb.UpdatePostRequest) (*pb.Post, error) {
	ctx := context.Background()
	res, err := c.client.UpdatePost(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *BlogClient) DeletePost(req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	ctx := context.Background()
	res, err := c.client.DeletePost(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
