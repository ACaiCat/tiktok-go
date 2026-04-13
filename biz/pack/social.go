package pack

import (
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/social"
)

func RespFollow(c *app.RequestContext) {
	c.JSON(http.StatusOK, social.FollowResp{
		Base: SuccessBase,
	})
}

func RespListFollower(c *app.RequestContext, userList []*model.SocialUser, total int) {
	c.JSON(http.StatusOK, social.ListFollowerResp{
		Base: SuccessBase,
		Data: &social.SocialUserData{
			Items: userList,
			Total: int32(total),
		},
	})
}

func RespListFollowing(c *app.RequestContext, userList []*model.SocialUser, total int) {
	c.JSON(http.StatusOK, social.ListFollowingResp{
		Base: SuccessBase,
		Data: &social.SocialUserData{
			Items: userList,
			Total: int32(total),
		},
	})
}

func RespListFriend(c *app.RequestContext, userList []*model.SocialUser, total int) {
	c.JSON(http.StatusOK, social.ListFriendResp{
		Base: SuccessBase,
		Data: &social.SocialUserData{
			Items: userList,
			Total: int32(total),
		},
	})
}
