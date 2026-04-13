package pack

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/video"
)

func RespFeedList(c *app.RequestContext, videoList []*model.Video) {
	c.JSON(consts.StatusOK, video.FeedResp{
		Base: SuccessBase,
		Data: &video.VideoList{
			Items: videoList,
		},
	})
}

func RespPopularList(c *app.RequestContext, videoList []*model.Video) {
	c.JSON(consts.StatusOK, video.PopularResp{
		Base: SuccessBase,
		Data: &video.VideoList{
			Items: videoList,
		},
	})
}

func RespVideoList(c *app.RequestContext, videoList []*model.Video, total int64) {
	c.JSON(consts.StatusOK, video.ListResp{
		Base: SuccessBase,
		Data: &video.VideoListWithTotal{
			Items: videoList,
			Total: total,
		},
	})
}

func RespPublish(c *app.RequestContext) {
	c.JSON(consts.StatusOK, video.PublishResp{
		Base: SuccessBase,
	})
}

func RespSearch(c *app.RequestContext, videoList []*model.Video) {
	c.JSON(consts.StatusOK, video.SearchResp{
		Base: SuccessBase,
		Data: &video.VideoList{
			Items: videoList,
		},
	})
}

func RespVisitVideo(c *app.RequestContext) {
	c.JSON(consts.StatusOK, video.VisitVideoResp{
		Base: SuccessBase,
	})
}
