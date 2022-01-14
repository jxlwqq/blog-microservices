package post

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	v1 "github.com/jxlwqq/blog-microservices/api/protobuf/post/v1"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/dbcontext"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	p1 := &v1.Post{
		Uuid:    uuid.NewString(),
		Title:   "Hello World",
		Content: "Hello World",
	}

	p2 := &v1.Post{
		Uuid:    uuid.NewString(),
		Title:   "Hello Go",
		Content: "Hello Go",
	}

	logger := log.New()
	path := config.GetPath()
	conf, err := config.Load(path)
	require.NoError(t, err)
	db, err := dbcontext.NewPostDB(conf)
	require.NoError(t, err)
	repo := NewRepository(logger, db)
	require.NotNil(t, repo)
	s := NewServer(logger, repo)
	require.NotNil(t, s)

	// Test CreatePost
	createResp, err := s.CreatePost(context.Background(), &v1.CreatePostRequest{Post: p1})
	require.NoError(t, err)
	require.NotNil(t, createResp.GetPost().GetId())
	createResp2, err := s.CreatePost(context.Background(), &v1.CreatePostRequest{Post: p2})
	require.NoError(t, err)
	require.NotNil(t, createResp2.GetPost().GetId())

	// Test GetPost
	getResp, err := s.GetPost(context.Background(), &v1.GetPostRequest{Id: createResp.GetPost().GetId()})
	fmt.Println(getResp.GetPost().GetUuid())
	require.NoError(t, err)
	require.NotNil(t, getResp.GetPost().GetTitle())
	require.Equal(t, p1.Uuid, getResp.GetPost().GetUuid())

	// Test UpdatePost
	p3 := getResp.GetPost()
	p3.Title = "Hello World2"
	updateResp, err := s.UpdatePost(context.Background(), &v1.UpdatePostRequest{Post: p3})
	require.NoError(t, err)
	require.True(t, updateResp.GetSuccess())

	// Test IncrementCommentCount and IncrementCommentsCountCompensate
	incrementResp, err := s.IncrementCommentsCount(context.Background(), &v1.IncrementCommentsCountRequest{Id: createResp.GetPost().GetId()})
	require.NoError(t, err)
	require.True(t, incrementResp.GetSuccess())
	getResp, err = s.GetPost(context.Background(), &v1.GetPostRequest{Id: createResp.GetPost().GetId()})
	require.NoError(t, err)
	require.Equal(t, uint32(1), getResp.GetPost().GetCommentsCount())
	incrementCompensateResp, err := s.IncrementCommentsCountCompensate(context.Background(), &v1.IncrementCommentsCountRequest{Id: createResp.GetPost().GetId()})
	require.NoError(t, err)
	require.True(t, incrementCompensateResp.GetSuccess())
	getResp, err = s.GetPost(context.Background(), &v1.GetPostRequest{Id: createResp.GetPost().GetId()})
	require.NoError(t, err)
	require.Equal(t, uint32(0), getResp.GetPost().GetCommentsCount())

	// Test DecrementCommentCount and DecrementCommentsCountCompensate
	incrementResp, err = s.IncrementCommentsCount(context.Background(), &v1.IncrementCommentsCountRequest{Id: createResp.GetPost().GetId()})
	require.NoError(t, err)
	require.True(t, incrementResp.GetSuccess())
	decrementResp, err := s.DecrementCommentsCount(context.Background(), &v1.DecrementCommentsCountRequest{Id: createResp.GetPost().GetId()})
	require.NoError(t, err)
	require.True(t, decrementResp.GetSuccess())
	getResp, err = s.GetPost(context.Background(), &v1.GetPostRequest{Id: createResp.GetPost().GetId()})
	require.NoError(t, err)
	require.Equal(t, uint32(0), getResp.GetPost().GetCommentsCount())
	decrementCompensateResp, err := s.DecrementCommentsCountCompensate(context.Background(), &v1.DecrementCommentsCountRequest{Id: createResp.GetPost().GetId()})
	require.NoError(t, err)
	require.True(t, decrementCompensateResp.GetSuccess())
	getResp, err = s.GetPost(context.Background(), &v1.GetPostRequest{Id: createResp.GetPost().GetId()})
	require.NoError(t, err)
	require.Equal(t, uint32(1), getResp.GetPost().GetCommentsCount())

	// Test ListPosts
	listResp, err := s.ListPosts(context.Background(), &v1.ListPostsRequest{Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.Equal(t, 2, len(listResp.GetPosts()))
	require.Equal(t, p1.Uuid, listResp.GetPosts()[0].GetUuid())
	require.Equal(t, p2.Uuid, listResp.GetPosts()[1].GetUuid())

	// Test DeletePost
	deleteResp, err := s.DeletePost(context.Background(), &v1.DeletePostRequest{Id: createResp.GetPost().GetId()})
	require.NoError(t, err)
	require.True(t, deleteResp.GetSuccess())
	deleteResp2, err := s.DeletePost(context.Background(), &v1.DeletePostRequest{Id: createResp2.GetPost().GetId()})
	require.NoError(t, err)
	require.True(t, deleteResp2.GetSuccess())

	// Test DeletePostCompensate
	deleteResp, err = s.DeletePostCompensate(context.Background(), &v1.DeletePostRequest{Id: createResp.GetPost().GetId()})
	require.NoError(t, err)
	require.True(t, deleteResp.GetSuccess())
	getResp, err = s.GetPost(context.Background(), &v1.GetPostRequest{Id: createResp.GetPost().GetId()})
	require.NoError(t, err)
	require.NotNil(t, getResp.GetPost().GetId())
}
