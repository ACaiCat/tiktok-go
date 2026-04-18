package main

import (
	"context"

	video "github.com/ACaiCat/tiktok-go/kitex_gen/video"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

// Feed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.FeedReq) (resp *video.FeedResp, err error) {
	// TODO: Your code here...
	return
}

// Publish implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Publish(ctx context.Context, req *video.PublishReq) (resp *video.PublishResp, err error) {
	// TODO: Your code here...
	return
}

// List implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) List(ctx context.Context, req *video.ListReq) (resp *video.ListResp, err error) {
	// TODO: Your code here...
	return
}

// Popular implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Popular(ctx context.Context, req *video.PopularReq) (resp *video.PopularResp, err error) {
	// TODO: Your code here...
	return
}

// Search implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Search(ctx context.Context, req *video.SearchReq) (resp *video.SearchResp, err error) {
	// TODO: Your code here...
	return
}

// VisitVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) VisitVideo(ctx context.Context, req *video.VisitVideoReq) (resp *video.VisitVideoResp, err error) {
	// TODO: Your code here...
	return
}
