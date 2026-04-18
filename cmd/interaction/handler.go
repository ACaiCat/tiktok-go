package main

import (
	"context"

	interaction "github.com/ACaiCat/tiktok-go/kitex_gen/interaction"
)

// InteractionServiceImpl implements the last service interface defined in the IDL.
type InteractionServiceImpl struct{}

// Like implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) Like(ctx context.Context, req *interaction.LikeReq) (resp *interaction.LikeResp, err error) {
	// TODO: Your code here...
	return
}

// ListLike implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) ListLike(ctx context.Context, req *interaction.ListLikeReq) (resp *interaction.ListLikeResp, err error) {
	// TODO: Your code here...
	return
}

// Comment implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) Comment(ctx context.Context, req *interaction.CommentReq) (resp *interaction.CommentResp, err error) {
	// TODO: Your code here...
	return
}

// ListComment implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) ListComment(ctx context.Context, req *interaction.ListCommentReq) (resp *interaction.ListCommentResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteComment implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) DeleteComment(ctx context.Context, req *interaction.DeleteCommentReq) (resp *interaction.DeleteCommentResp, err error) {
	// TODO: Your code here...
	return
}
