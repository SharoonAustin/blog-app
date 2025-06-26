package service

import (
	pb "blog-app/proto"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
	"net/http"
	"sync"
	"time"
)

const (
	postNotFound = "couldn't find post for the given post id %s"
)

type BlogServer struct {
	pb.UnimplementedPostServiceServer
	log     *logrus.Logger
	postMap map[string]*pb.Post
	mu      sync.RWMutex
}

func NewBlogServer() *BlogServer {
	return &BlogServer{
		log:     logrus.New(),
		postMap: make(map[string]*pb.Post),
	}
}

func (b *BlogServer) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	b.log.Info("invoked create post", req)

	if req.GetTitle() == "" {
		return nil, errors.New("title field is required")
	}
	if req.GetContent() == "" {
		return nil, errors.New("content field is required")
	}
	if req.GetAuthor() == "" {
		return nil, errors.New("author field is required")
	}
	if req.GetPublicationDate() == "" {
		return nil, errors.New("publication_date field is required")
	}

	_, err := time.Parse("2006-01-02", req.GetPublicationDate())
	if err != nil {
		return nil, errors.New("publication_date must be in YYYY-MM-DD format")
	}

	for _, post := range b.postMap {
		if post.GetAuthor() == req.GetAuthor() && post.GetTitle() == req.GetTitle() {
			return nil, errors.New("a post with the same title by this author already exists")
		}
	}

	post := &pb.Post{
		PostId:          uuid.NewString(),
		Title:           req.GetTitle(),
		Content:         req.GetContent(),
		Author:          req.GetAuthor(),
		PublicationDate: req.GetPublicationDate(),
		Tags:            req.GetTags(),
	}

	b.mu.Lock()
	b.postMap[post.GetPostId()] = post
	defer b.mu.Unlock()

	b.log.Info("post is created successfully", post)
	return post, nil
}

func (b *BlogServer) ReadPost(ctx context.Context, req *pb.ReadPostRequest) (*pb.Post, error) {
	b.log.Info("invoked read post", req)
	b.mu.RLock()
	post, isPresent := b.postMap[req.GetPostId()]
	b.mu.RUnlock()

	if !isPresent {
		b.log.Error(fmt.Sprintf(postNotFound, req.PostId))
		return nil, status.Error(http.StatusNotFound, fmt.Sprintf(postNotFound, req.PostId))
	}
	b.log.Info("post is retrieved successfully", post)
	return post, nil
}

func (b *BlogServer) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.Post, error) {
	b.log.Info("invoked update post", req)
	postInfo, err := b.ReadPost(ctx, &pb.ReadPostRequest{PostId: req.PostId})
	if err != nil {
		b.log.Error("error occurred while updating the post", err)
		return nil, err
	}
	if req.Title == "" {
		req.Title = postInfo.Title
	}
	if req.Content == "" {
		req.Content = postInfo.Content
	}

	if req.Author == "" {
		req.Author = postInfo.Author
	}

	req.Tags = append(req.Tags, postInfo.Tags...)

	post := &pb.Post{
		PostId:          req.GetPostId(),
		Title:           req.GetTitle(),
		Content:         req.GetContent(),
		Author:          req.GetAuthor(),
		PublicationDate: postInfo.PublicationDate,
		Tags:            req.GetTags(),
	}
	b.mu.Lock()
	b.postMap[req.PostId] = post
	b.mu.Unlock()
	b.log.Info("post is retrieved successfully", post)
	return post, nil
}

func (b *BlogServer) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	b.log.Info("invoked delete post", req)
	_, err := b.ReadPost(ctx, &pb.ReadPostRequest{PostId: req.PostId})
	if err != nil {
		b.log.Error("error occurred while deleting the post", err)
		return nil, err
	}
	b.mu.Lock()
	delete(b.postMap, req.PostId)
	b.mu.Unlock()
	b.log.Info("post is deleted successfully")
	return &pb.DeletePostResponse{Message: "successfully deleted the post"}, nil
}
