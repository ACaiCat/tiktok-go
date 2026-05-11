package jwt

import (
	"os"
	"testing"
	"time"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/config"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

const (
	testAccessSecret  = "114514"
	testRefreshSecret = "1919810"
	testUserID        = int64(1)
	testAccessToken   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ0b2tlbl90eXBlIjowLCJpc3MiOiJDYWkiLCJle" +
		"HAiOjE3Nzg1MDYxODUsImlhdCI6MTc3ODQ5ODk4NX0.gQMvOj1Fh-Z8SSuc4qAOcMpn4w_1T2fgFz1Nvbk2upc"
	testRefreshToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ0b2tlbl90eXBlIjoxLCJpc3MiOiJDYWkiLCJleH" +
		"AiOjE3NzkxMDM3ODUsImlhdCI6MTc3ODQ5ODk4NX0.bUbq8s7s3XgATxRaX0uaPLX9rZBGMI8leo6cSrg9Npo"
	testTokenIssueAt = 1778498985000
)

func TestMain(m *testing.M) {
	config.AppConfig.JWT.AccessSecret = testAccessSecret
	config.AppConfig.JWT.RefreshSecret = testRefreshSecret
	code := m.Run()
	os.Exit(code)
}

func TestCreateToken(t *testing.T) {
	type testCase struct {
		tokenType int8
		userID    int64
		timestamp int64
		wantToken string
		wantErr   bool
	}

	testCases := map[string]testCase{
		"validate access token": {
			tokenType: constants.TypeAccessToken,
			userID:    testUserID,
			timestamp: testTokenIssueAt,
			wantToken: testAccessToken,
			wantErr:   false,
		},
		"validate refresh token": {
			tokenType: constants.TypeRefreshToken,
			userID:    testUserID,
			timestamp: testTokenIssueAt,
			wantToken: testRefreshToken,
			wantErr:   false,
		},
		"invalid token type": {
			tokenType: 3,
			userID:    testUserID,
			timestamp: testTokenIssueAt,
			wantToken: testRefreshToken,
			wantErr:   true,
		},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			testTime := time.UnixMilli(tc.timestamp)
			Mock(time.Now).Return(testTime).Build()

			token, err := CreateToken(tc.tokenType, tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tc.wantToken, token)
		})
	}
}

func TestVerifyToken(t *testing.T) {
	type testCase struct {
		token     string
		tokenType int8
		userID    int64
		timestamp int64
		wantErr   error
	}

	testCases := map[string]testCase{
		"validate access token": {
			token:     testAccessToken,
			tokenType: constants.TypeAccessToken,
			userID:    testUserID,
			timestamp: testTokenIssueAt,
			wantErr:   nil,
		},
		"validate refresh token": {
			token:     testRefreshToken,
			tokenType: constants.TypeRefreshToken,
			userID:    testUserID,
			timestamp: testTokenIssueAt,
			wantErr:   nil,
		},
		"expired token access token": {
			token:     testAccessToken,
			tokenType: constants.TypeAccessToken,
			userID:    testUserID,
			timestamp: testTokenIssueAt + constants.AccessTokenExpiration.Milliseconds(),
			wantErr:   errno.AuthAccessExpiredErr,
		},
		"expired token refresh token": {
			token:     testRefreshToken,
			tokenType: constants.TypeRefreshToken,
			userID:    testUserID,
			timestamp: testTokenIssueAt + constants.RefreshTokenExpiration.Milliseconds(),
			wantErr:   errno.AuthRefreshExpiredErr,
		},
		"missing token": {
			token:     "",
			tokenType: constants.TypeAccessToken,
			userID:    testUserID,
			timestamp: testTokenIssueAt,
			wantErr:   errno.AuthMissingErr,
		},
		"unmatched token type": {
			token:     testAccessToken,
			tokenType: constants.TypeRefreshToken,
			userID:    testUserID,
			timestamp: testTokenIssueAt,
			wantErr:   errno.AuthErr.WithMessage("令牌类型不匹配"),
		},
		"invalid token": {
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ0b2tlbl90eXBlIjoxLCJpc3MiOiJDYWkiLCJleHA" +
				"iOjE3NzkxMDM3ODUsImlhdCI6MTc3ODQ5ODk4NX0.RBus5dO_3J_IHa-dT7_pN6H8uXp4GFFRLlCR9cBMACI",
			tokenType: constants.TypeRefreshToken,
			userID:    testUserID,
			timestamp: testTokenIssueAt,
			wantErr:   errno.AuthErr.WithMessage("令牌无效"),
		},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			testTime := time.UnixMilli(tc.timestamp)
			Mock(time.Now).Return(testTime).Build()

			userID, err := ValidateToken(tc.token, tc.tokenType)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			assert.Equal(t, tc.userID, userID)
		})
	}
}
