package pack

import (
	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/social"
	"github.com/cloudwego/hertz/pkg/app"
)

func RespFollow(c *app.RequestContext) {
	c.JSON(200, social.FollowResp{
		Base: SuccessBase,
	})
}

func RespListFollower(c *app.RequestContext, userList []*model.SocialUser, total int) {
	c.JSON(200, social.ListFollowerResp{
		Base: SuccessBase,
		Data: &social.SocialUserData{
			Items: userList,
			Total: int32(total),
		},
	})
}

func RespListFollowing(c *app.RequestContext, userList []*model.SocialUser, total int) {
	c.JSON(200, social.ListFollowingResp{
		Base: SuccessBase,
		Data: &social.SocialUserData{
			Items: userList,
			Total: int32(total),
		},
	})
}

func RespListFriend(c *app.RequestContext, userList []*model.SocialUser, total int) {
	c.JSON(200, social.ListFriendResp{
		Base: SuccessBase,
		Data: &social.SocialUserData{
			Items: userList,
			Total: int32(total),
		},
	})
}
