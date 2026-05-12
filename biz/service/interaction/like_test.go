package service

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/biz/model/interaction"
)

func TestInteractionService_LikeVideo(t *testing.T) {
	type testCase struct {
		req            *interaction.LikeReq
		mockCommentErr error
		expectError    string
	}

	testCases := map[string]testCase{
		"missing target": {
			req:         &interaction.LikeReq{ActionType: interaction.LikeActionType_ADD},
			expectError: "视频ID或评论ID不能为空",
		},
		"both targets": {
			req:         &interaction.LikeReq{VideoID: new("1"), CommentID: new("2"), ActionType: interaction.LikeActionType_ADD},
			expectError: "视频ID和评论ID不能同时存在",
		},
		"comment success": {
			req: &interaction.LikeReq{CommentID: new("2"), ActionType: interaction.LikeActionType_ADD},
		},
		"comment error": {
			req:            &interaction.LikeReq{CommentID: new("2"), ActionType: interaction.LikeActionType_ADD},
			mockCommentErr: assert.AnError,
			expectError:    assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*InteractionService).likeCommentByID).To(
				func(_ *InteractionService, commentIDStr string, userID int64, actionType interaction.LikeActionType) error {
					return tc.mockCommentErr
				}).Build()

			mockey.Mock(NewInteractionService).To(func(_ context.Context) *InteractionService {
				return &InteractionService{}
			}).Build()

			err := NewInteractionService(context.Background()).LikeVideo(tc.req, 1)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
		})
	}
}
