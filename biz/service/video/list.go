package service

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/biz/model/interaction"
	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/video"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/ACaiCat/tiktok-go/pkg/utils"
)

func (s *VideoService) GetVideoList(req *video.ListReq) ([]*model.Video, int64, error) {
	pageSize, pageNum := utils.NormalizePage(req.PageSize, req.PageNum, constants.DefaultVideoPageSize, constants.MaxVideoPageSize)

	userID, err := strconv.ParseInt(req.UserID, 10, 64)

	if err != nil {
		return nil, 0, errno.ParamErr.WithError(err)
	}

	var videosDao []*modelDao.Video
	var total int64

	videoIDs, cachedTotal, err := s.videoCache.GetUserVideoList(s.ctx, userID, pageSize, pageNum)
	if err == nil {
		videosDao, err := s.getVideosByIDs(videoIDs)
		if err != nil {
			return nil, 0, err
		}
		return VideosDaoToDto(videosDao), cachedTotal, nil
	}
	if !errors.Is(err, redis.Nil) {
		hlog.CtxErrorf(s.ctx, "service.GetVideoList cache read failed: %v", err)
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		exist, err := s.userDao.WithTx(tx).IsUserExists(s.ctx, userID)
		if err != nil {
			return errors.WithMessagef(err, "service.GetVideoList check user exists userID: %d", userID)
		}
		if !exist {
			return errno.UserIsNotExistErr
		}

		videosDao, err = s.videoDao.WithTx(tx).GetVideosByUserID(s.ctx, userID, pageSize, pageNum)
		if err != nil {
			return errors.WithMessagef(err, "service.GetVideoList: db.GetVideosByUserID failed, userID=%d", userID)
		}

		total, err = s.videoDao.WithTx(tx).GetVideoCountByUserID(s.ctx, userID)
		if err != nil {
			return errors.WithMessagef(err, "service.GetVideoList: db.GetVideoCountByUserID failed, userID=%d", userID)
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	go func() {
		if err := s.videoCache.SetUserVideoList(context.Background(), userID, pageSize, pageNum, total, videosDao); err != nil {
			hlog.Errorf("service.GetVideoList cache list write failed: %v", err)
		}
		if len(videosDao) > 0 {
			if err := s.videoCache.SetVideos(context.Background(), videosDao); err != nil {
				hlog.Errorf("service.GetVideoList cache detail write failed: %v", err)
			}
		}
	}()

	return VideosDaoToDto(videosDao), total, nil
}

func (s *VideoService) GetLikedVideos(req *interaction.ListLikeReq) ([]*model.Video, error) {
	pageSize, pageNum := utils.NormalizePage(req.PageSize, req.PageNum, constants.DefaultVideoPageSize, constants.MaxVideoPageSize)

	userID, err := strconv.ParseInt(req.UserID, 10, 64)
	if err != nil {
		return nil, errno.ParamErr.WithError(err)
	}

	likedVideoIDs, err := s.userCache.GetLikedVideos(s.ctx, userID)
	if err == nil && len(likedVideoIDs) > 0 {
		start := pageSize * pageNum
		if start >= len(likedVideoIDs) {
			return []*model.Video{}, nil
		}
		end := start + pageSize
		if end > len(likedVideoIDs) {
			end = len(likedVideoIDs)
		}

		likedVideos, err := s.getVideosByIDs(likedVideoIDs[start:end])
		if err != nil {
			return nil, err
		}
		return VideosDaoToDto(likedVideos), nil
	}
	if err != nil && !errors.Is(err, redis.Nil) {
		hlog.CtxErrorf(s.ctx, "service.GetLikedVideos cache read failed: %v", err)
	}

	likedVideos, err := s.videoDao.GetUserLikeList(s.ctx, userID, pageSize, pageNum)
	if err != nil {
		return nil, errors.WithMessagef(err, "service.GetLikedVideos: db.GetUserLikeList failed, userID=%d", userID)
	}

	if len(likedVideos) > 0 {
		go func() {
			if err := s.videoCache.SetVideos(context.Background(), likedVideos); err != nil {
				hlog.Errorf("service.GetLikedVideos cache detail write failed: %v", err)
			}
		}()
	}

	videos := VideosDaoToDto(likedVideos)

	return videos, nil
}
