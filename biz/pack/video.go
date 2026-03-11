package pack

import (
	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/video"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func RespFeedList(c *app.RequestContext, videoList []*model.Video) {
	c.JSON(consts.StatusOK, video.FeedResp{
		Base:  SuccessBase,
		Items: videoList,
	})
}

func RespPopularList(c *app.RequestContext, videoList []*model.Video) {
	c.JSON(consts.StatusOK, video.PopularResp{
		Base:  SuccessBase,
		Items: videoList,
	})
}

func RespAuthorList(c *app.RequestContext, videoList []*model.Video) {
	c.JSON(consts.StatusOK, video.ListResp{
		Base:  SuccessBase,
		Items: videoList,
	})
}

func RespPublish(c *app.RequestContext) {
	c.JSON(consts.StatusOK, video.PublishResp{
		Base: SuccessBase,
	})
}
