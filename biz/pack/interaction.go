package pack

import (
	"github.com/ACaiCat/tiktok-go/biz/model/interaction"
	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func RespLike(c *app.RequestContext) {
	c.JSON(consts.StatusOK, interaction.LikeResp{
		Base: SuccessBase,
	})
}

func RespLikeList(c *app.RequestContext, videoList []*model.Video) {
	c.JSON(consts.StatusOK, interaction.ListLikeResp{
		Base: SuccessBase,
		Data: &interaction.VideoData{
			Items: videoList,
		},
	})
}

func RespComment(c *app.RequestContext) {
	c.JSON(consts.StatusOK, interaction.CommentResp{
		Base: SuccessBase,
	})
}

func RespListComment(c *app.RequestContext, comments []*model.Comment) {
	c.JSON(consts.StatusOK, interaction.ListCommentResp{
		Base: SuccessBase,
		Data: &interaction.CommentData{
			Items: comments,
		},
	})
}

func RespDeleteComment(c *app.RequestContext) {
	c.JSON(consts.StatusOK, interaction.DeleteCommentResp{
		Base: SuccessBase,
	})
}
