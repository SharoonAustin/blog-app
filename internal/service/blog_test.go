package service

import (
	"context"
	"fmt"
	"sync"
	"testing"

	pb "blog-app/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBlogServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BlogServer Suite")
}

var _ = Describe("BlogServer", func() {
	var (
		server *BlogServer
		ctx    context.Context
	)

	BeforeEach(func() {
		server = NewBlogServer()
		ctx = context.Background()
	})

	Describe("CreatePost", func() {
		It("should create a new post successfully", func() {
			req := &pb.CreatePostRequest{
				Title:           "Test Title",
				Content:         "Test Content",
				Author:          "Author1",
				PublicationDate: "2025-06-26",
				Tags:            []string{"tag1", "tag2"},
			}

			post, err := server.CreatePost(ctx, req)
			Expect(err).ToNot(HaveOccurred())
			Expect(post).ToNot(BeNil())
			Expect(post.PostId).ToNot(BeEmpty())
			Expect(post.Title).To(Equal(req.Title))
			Expect(post.Author).To(Equal(req.Author))
		})

		It("should fail when title is missing", func() {
			req := &pb.CreatePostRequest{
				Content:         "Content",
				Author:          "Author1",
				PublicationDate: "2025-06-26",
				Tags:            []string{"tag1"},
			}
			_, err := server.CreatePost(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("title field is required"))
		})

		It("should fail when publication date is missing", func() {
			req := &pb.CreatePostRequest{
				Title:   "Title",
				Content: "Content",
				Author:  "Author1",
				Tags:    []string{"tag1"},
			}
			_, err := server.CreatePost(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("publication_date field is required"))
		})

		It("should fail when publication date format is invalid", func() {
			req := &pb.CreatePostRequest{
				Title:           "Title",
				Content:         "Content",
				Author:          "Author1",
				PublicationDate: "26-06-2025",
				Tags:            []string{"tag1"},
			}
			_, err := server.CreatePost(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("publication_date must be in YYYY-MM-DD format"))
		})

		It("should fail when a post with same title by same author exists", func() {
			req1 := &pb.CreatePostRequest{
				Title:           "Duplicate Title",
				Content:         "Content 1",
				Author:          "Author1",
				PublicationDate: "2025-06-26",
				Tags:            []string{"tag1"},
			}
			_, err := server.CreatePost(ctx, req1)
			Expect(err).ToNot(HaveOccurred())

			req2 := &pb.CreatePostRequest{
				Title:           "Duplicate Title",
				Content:         "Content 2",
				Author:          "Author1",
				PublicationDate: "2025-06-27",
				Tags:            []string{"tag2"},
			}
			_, err = server.CreatePost(ctx, req2)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("a post with the same title by this author already exists"))
		})
	})

	Describe("ReadPost", func() {
		It("should retrieve an existing post", func() {
			req := &pb.CreatePostRequest{
				Title:           "ReadTest",
				Content:         "Some content",
				Author:          "Author2",
				PublicationDate: "2025-06-26",
				Tags:            []string{"tagA"},
			}
			createdPost, _ := server.CreatePost(ctx, req)

			readReq := &pb.ReadPostRequest{PostId: createdPost.PostId}
			post, err := server.ReadPost(ctx, readReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(post.PostId).To(Equal(createdPost.PostId))
			Expect(post.Title).To(Equal(createdPost.Title))
		})

		It("should return error if post does not exist", func() {
			readReq := &pb.ReadPostRequest{PostId: "non-existent-id"}
			post, err := server.ReadPost(ctx, readReq)
			Expect(post).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("UpdatePost", func() {
		It("should update an existing post", func() {
			createReq := &pb.CreatePostRequest{
				Title:           "Old Title",
				Content:         "Old Content",
				Author:          "Author3",
				PublicationDate: "2025-06-26",
				Tags:            []string{"tagOld"},
			}
			createdPost, _ := server.CreatePost(ctx, createReq)

			updateReq := &pb.UpdatePostRequest{
				PostId:  createdPost.PostId,
				Title:   "New Title",
				Content: "New Content",
				Author:  "Author3Updated",
				Tags:    []string{"tagNew"},
			}

			updatedPost, err := server.UpdatePost(ctx, updateReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(updatedPost.Title).To(Equal("New Title"))
			Expect(updatedPost.Content).To(Equal("New Content"))
			Expect(updatedPost.Author).To(Equal("Author3Updated"))
			Expect(updatedPost.Tags).To(ContainElement("tagNew"))
			Expect(updatedPost.PublicationDate).To(Equal(createdPost.PublicationDate))
		})

		It("should return error when updating non-existent post", func() {
			updateReq := &pb.UpdatePostRequest{
				PostId:  "non-existent-id",
				Title:   "No Title",
				Content: "No Content",
				Author:  "Nobody",
				Tags:    []string{},
			}

			post, err := server.UpdatePost(ctx, updateReq)
			Expect(post).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("DeletePost", func() {
		It("should delete an existing post", func() {
			createReq := &pb.CreatePostRequest{
				Title:           "To be deleted",
				Content:         "Content to delete",
				Author:          "AuthorDelete",
				PublicationDate: "2025-06-26",
				Tags:            []string{"tagDel"},
			}
			createdPost, _ := server.CreatePost(ctx, createReq)

			delReq := &pb.DeletePostRequest{PostId: createdPost.PostId}
			resp, err := server.DeletePost(ctx, delReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Message).To(Equal("successfully deleted the post"))

			_, err = server.ReadPost(ctx, &pb.ReadPostRequest{PostId: createdPost.PostId})
			Expect(err).To(HaveOccurred())
		})

		It("should return error when deleting non-existent post", func() {
			delReq := &pb.DeletePostRequest{PostId: "non-existent-id"}
			resp, err := server.DeletePost(ctx, delReq)
			Expect(resp).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Concurrency", func() {
		It("should handle concurrent reads and writes safely", func() {
			wg := sync.WaitGroup{}
			numOps := 100

			createReq := &pb.CreatePostRequest{
				Title:           "Concurrent",
				Content:         "Testing concurrency",
				Author:          "ConcurrentAuthor",
				PublicationDate: "2025-06-26",
				Tags:            []string{"concurrent"},
			}
			createdPost, _ := server.CreatePost(ctx, createReq)

			for i := 0; i < numOps; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, err := server.ReadPost(ctx, &pb.ReadPostRequest{PostId: createdPost.PostId})
					Expect(err).ToNot(HaveOccurred())
				}()
			}

			for i := 0; i < numOps; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					//ch <- i
					_, err := server.UpdatePost(ctx, &pb.UpdatePostRequest{
						PostId:  createdPost.PostId,
						Title:   fmt.Sprintf("Title %d", i),
						Content: "Updated content",
						Author:  "ConcurrentAuthor",
						Tags:    []string{"updated"},
					})
					Expect(err).ToNot(HaveOccurred())
				}(i)
			}

			wg.Wait()

			finalPost, err := server.ReadPost(ctx, &pb.ReadPostRequest{PostId: createdPost.PostId})
			Expect(err).ToNot(HaveOccurred())
			Expect(finalPost.Title).To(ContainSubstring("Title"))
		})

	})
})
