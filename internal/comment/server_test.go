package comment

import (
	"context"
	"github.com/google/uuid"
	v1 "github.com/jxlwqq/blog-microservices/api/protobuf/comment/v1"
	"github.com/jxlwqq/blog-microservices/internal/pkg/config"
	"github.com/jxlwqq/blog-microservices/internal/pkg/dbcontext"
	"github.com/jxlwqq/blog-microservices/internal/pkg/log"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServer(t *testing.T) {
	c1 := &v1.Comment{
		Uuid:    uuid.NewString(),
		PostId:  1,
		UserId:  1,
		Content: "Hello World",
	}

	c2 := &v1.Comment{
		Uuid:    uuid.NewString(),
		PostId:  1,
		UserId:  2,
		Content: "Hello Go",
	}

	c3 := &v1.Comment{
		Uuid:    uuid.NewString(),
		PostId:  1,
		UserId:  3,
		Content: "Hello Python",
	}

	logger := log.New()
	path := config.GetPath()
	conf, err := config.Load(path)
	require.NoError(t, err)
	db, err := dbcontext.NewCommentDB(conf)
	require.NoError(t, err)
	repo := NewRepository(logger, db)
	require.NotNil(t, repo)
	s := NewServer(logger, repo)
	require.NotNil(t, s)

	// Test Create
	createResp1, err := s.CreateComment(context.Background(), &v1.CreateCommentRequest{Comment: c1})
	require.NoError(t, err)
	require.NotNil(t, createResp1)
	require.Equal(t, c1.Uuid, createResp1.Comment.Uuid)

	createResp2, err := s.CreateComment(context.Background(), &v1.CreateCommentRequest{Comment: c2})
	require.NoError(t, err)
	require.NotNil(t, createResp2)
	require.Equal(t, c2.Uuid, createResp2.Comment.Uuid)

	createResp3, err := s.CreateComment(context.Background(), &v1.CreateCommentRequest{Comment: c3})
	require.NoError(t, err)
	require.NotNil(t, createResp3)
	require.Equal(t, c3.Uuid, createResp3.Comment.Uuid)

	// Test CreateCompensate
	createCompensateResp3, err := s.CreateCommentCompensate(context.Background(), &v1.CreateCommentRequest{Comment: c3})
	require.NoError(t, err)
	require.NotNil(t, createCompensateResp3)
	require.Nil(t, createCompensateResp3.GetComment())

	// Test List
	listResp, err := s.ListCommentsByPostID(context.Background(), &v1.ListCommentsByPostIDRequest{PostId: uint64(1)})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Equal(t, 2, len(listResp.Comments))
	require.Equal(t, c1.Uuid, listResp.Comments[0].Uuid)
	require.Equal(t, c2.Uuid, listResp.Comments[1].Uuid)

	// Test Get
	getResp1, err := s.GetComment(context.Background(), &v1.GetCommentRequest{Id: createResp1.GetComment().GetId()})
	require.NoError(t, err)
	require.NotNil(t, getResp1)
	require.Equal(t, c1.Uuid, getResp1.Comment.Uuid)

	getResp2, err := s.GetCommentByUUID(context.Background(), &v1.GetCommentByUUIDRequest{Uuid: c2.Uuid})
	require.NoError(t, err)
	require.NotNil(t, getResp2)
	require.Equal(t, c2.Uuid, getResp2.Comment.Uuid)

	// Test Update
	c4 := &v1.Comment{
		Id:      createResp1.GetComment().GetId(),
		Content: "Hello World Updated",
	}
	updateResp1, err := s.UpdateComment(context.Background(), &v1.UpdateCommentRequest{Comment: c4})
	require.NoError(t, err)
	require.NotNil(t, updateResp1)
	require.True(t, updateResp1.Success)

	// Test Delete
	deleteResp1, err := s.DeleteComment(context.Background(), &v1.DeleteCommentRequest{Id: createResp1.GetComment().GetId()})
	require.NoError(t, err)
	require.NotNil(t, deleteResp1)

	deleteResp2, err := s.DeleteComment(context.Background(), &v1.DeleteCommentRequest{Id: createResp2.GetComment().GetId()})
	require.NoError(t, err)
	require.NotNil(t, deleteResp2)

	// Test DeleteCompensate
	deleteCompensateResp1, err := s.DeleteCommentCompensate(context.Background(), &v1.DeleteCommentRequest{Id: createResp1.GetComment().GetId()})
	require.NoError(t, err)
	require.NotNil(t, deleteCompensateResp1)
	getResp1, err = s.GetComment(context.Background(), &v1.GetCommentRequest{Id: createResp1.GetComment().GetId()})
	require.NoError(t, err)
	require.NotNil(t, getResp1)
	require.Equal(t, c1.Uuid, getResp1.Comment.Uuid)
	deleteResp1, err = s.DeleteComment(context.Background(), &v1.DeleteCommentRequest{Id: createResp1.GetComment().GetId()})
	require.NoError(t, err)
	require.NotNil(t, deleteResp1)

}
