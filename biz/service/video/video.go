package service

import (
	"github.com/ACaiCat/tiktok-go/pkg/cache"
	videoCache "github.com/ACaiCat/tiktok-go/pkg/cache/video"
	"github.com/ACaiCat/tiktok-go/pkg/db"
	commentDao "github.com/ACaiCat/tiktok-go/pkg/db/comment"
	likeDao "github.com/ACaiCat/tiktok-go/pkg/db/like"
	videoDao "github.com/ACaiCat/tiktok-go/pkg/db/vedio"
)

type VideoService struct {
	videoDao   *videoDao.VideoDao
	commentDao *commentDao.CommentDao
	likeDao    *likeDao.LikeDao
	cache      *videoCache.VideoCache
}

func NewVideoService() *VideoService {
	return &VideoService{
		videoDao:   videoDao.NewVideoDao(db.DB),
		commentDao: commentDao.NewCommentDao(db.DB),
		likeDao:    likeDao.NewLikeDao(db.DB),
		cache:      videoCache.NewVideoCache(cache.Cache),
	}
}
