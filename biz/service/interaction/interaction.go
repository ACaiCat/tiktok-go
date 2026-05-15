package service

import (
	"context"

	"github.com/ACaiCat/tiktok-go/pkg/cache"
	userCache "github.com/ACaiCat/tiktok-go/pkg/cache/user"
	videoCache "github.com/ACaiCat/tiktok-go/pkg/cache/video"
	"github.com/ACaiCat/tiktok-go/pkg/db"
	commentDao "github.com/ACaiCat/tiktok-go/pkg/db/comment"
	likeDao "github.com/ACaiCat/tiktok-go/pkg/db/like"
	videoDao "github.com/ACaiCat/tiktok-go/pkg/db/video"
)

type InteractionService struct {
	commentDao *commentDao.CommentDao
	videoDao   *videoDao.VideoDao
	likeDao    *likeDao.LikeDao
	userCache  *userCache.UserCache
	videoCache *videoCache.VideoCache
	ctx        context.Context
}

func NewInteractionService(ctx context.Context) *InteractionService {
	return &InteractionService{
		commentDao: commentDao.NewCommentDao(db.DB),
		likeDao:    likeDao.NewLikeDao(db.DB),
		videoDao:   videoDao.NewVideoDao(db.DB),
		userCache:  userCache.NewUserCache(cache.Cache),
		videoCache: videoCache.NewVideoCache(cache.Cache),
		ctx:        ctx,
	}
}
