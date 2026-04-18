package main

import (
	"context"

	social "github.com/ACaiCat/tiktok-go/kitex_gen/social"
)

// SocialServiceImpl implements the last service interface defined in the IDL.
type SocialServiceImpl struct{}

// Follow implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) Follow(ctx context.Context, req *social.FollowReq) (resp *social.FollowResp, err error) {
	// TODO: Your code here...
	return
}

// ListFollowing implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) ListFollowing(ctx context.Context, req *social.ListFollowingReq) (resp *social.ListFollowingResp, err error) {
	// TODO: Your code here...
	return
}

// ListFollower implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) ListFollower(ctx context.Context, req *social.ListFollowerReq) (resp *social.ListFollowerResp, err error) {
	// TODO: Your code here...
	return
}

// ListFriend implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) ListFriend(ctx context.Context, req *social.ListFriendReq) (resp *social.ListFriendResp, err error) {
	// TODO: Your code here...
	return
}
