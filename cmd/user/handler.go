package main

import (
	"context"

	user "github.com/ACaiCat/tiktok-go/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	// TODO: Your code here...
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	// TODO: Your code here...
	return
}

// Refresh implements the UserServiceImpl interface.
func (s *UserServiceImpl) Refresh(ctx context.Context, req *user.RefreshReq) (resp *user.RefreshResp, err error) {
	// TODO: Your code here...
	return
}

// Info implements the UserServiceImpl interface.
func (s *UserServiceImpl) Info(ctx context.Context, req *user.InfoReq) (resp *user.InfoResp, err error) {
	// TODO: Your code here...
	return
}

// UploadAvatar implements the UserServiceImpl interface.
func (s *UserServiceImpl) UploadAvatar(ctx context.Context, req *user.UploadAvatarReq) (resp *user.UploadAvatarResp, err error) {
	// TODO: Your code here...
	return
}

// MFAQRCode implements the UserServiceImpl interface.
func (s *UserServiceImpl) MFAQRCode(ctx context.Context, req *user.MFAQRCodeReq) (resp *user.MFAQRCodeResp, err error) {
	// TODO: Your code here...
	return
}

// BindMFA implements the UserServiceImpl interface.
func (s *UserServiceImpl) BindMFA(ctx context.Context, req *user.BindMFAReq) (resp *user.BindMFAResp, err error) {
	// TODO: Your code here...
	return
}

// SearchImage implements the UserServiceImpl interface.
func (s *UserServiceImpl) SearchImage(ctx context.Context, req *user.SearchImageReq) (resp *user.SearchImageResp, err error) {
	// TODO: Your code here...
	return
}
