package main

import (
	"blog-app/client"
	pb "blog-app/proto"
	"fmt"
	"log"
)

func main() {
	cli, err := client.NewBlogClient("localhost:50051")
	if err != nil {
		panic(err)
	}

	createReq := &pb.CreatePostRequest{
		Title:           "Harry Potter and the Deathly Hollows",
		Content:         "Fantasy",
		Author:          "J.K. Rowling",
		PublicationDate: "2006-05-05",
		Tags:            []string{"HP", "Harry Potter"},
	}

	post, err := cli.CreatePost(createReq)
	if err != nil {
		log.Fatalf(fmt.Sprintf("error occurred while creating post %v", err))
	}

	fmt.Println("Created Post : ", post)

	post, err = cli.GetPost(&pb.ReadPostRequest{PostId: post.PostId})
	if err != nil {
		log.Fatalf(fmt.Sprintf("error occurred while getting post %v", err))
	}

	fmt.Println("Get Post : ", post)

	post, err = cli.UpdatePost(&pb.UpdatePostRequest{PostId: post.PostId, Tags: []string{"Magic"}})
	if err != nil {
		log.Fatalf(fmt.Sprintf("error occurred while updating post %v", err))
	}

	fmt.Println("Update Post : ", post)

	delResp, err := cli.DeletePost(&pb.DeletePostRequest{PostId: post.PostId})
	if err != nil {
		log.Fatalf(fmt.Sprintf("error occurred while deleting post %v", err))
	}

	fmt.Println("Delete Post : ", delResp.Message)

}
