package errno

import (
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestNewErrNo(t *testing.T) {
	type testCase struct {
		code    int32
		msg     string
		wantErr ErrNo
	}

	testCases := map[string]testCase{
		"create success error": {
			code:    SuccessCode,
			msg:     "成功",
			wantErr: ErrNo{ErrCode: SuccessCode, ErrMsg: "成功"},
		},
		"create service error": {
			code:    ServiceErrCode,
			msg:     "服务内部错误",
			wantErr: ErrNo{ErrCode: ServiceErrCode, ErrMsg: "服务内部错误"},
		},
		"create param error": {
			code:    ParamErrCode,
			msg:     "参数错误",
			wantErr: ErrNo{ErrCode: ParamErrCode, ErrMsg: "参数错误"},
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			err := NewErrNo(tc.code, tc.msg)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestErrNo_Error(t *testing.T) {
	type testCase struct {
		err     ErrNo
		wantMsg string
	}

	testCases := map[string]testCase{
		"success error message": {
			err:     Success,
			wantMsg: "成功",
		},
		"service error message": {
			err:     ServiceErr,
			wantMsg: "服务器内部错误",
		},
		"auth error message": {
			err:     AuthErr,
			wantMsg: "认证失败",
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			assert.Equal(t, tc.wantMsg, tc.err.Error())
		})
	}
}

func TestErrNo_WithMessage(t *testing.T) {
	type testCase struct {
		base    ErrNo
		newMsg  string
		wantMsg string
	}

	testCases := map[string]testCase{
		"override service error message": {
			base:    ServiceErr,
			newMsg:  "自定义错误",
			wantMsg: "自定义错误",
		},
		"override avatar too large message": {
			base:    AvatarTooLargeErr,
			newMsg:  "图片超过最大限制，请上传小于5MB的图片",
			wantMsg: "图片超过最大限制，请上传小于5MB的图片",
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			result := tc.base.WithMessage(tc.newMsg)
			assert.Equal(t, tc.wantMsg, result.ErrMsg)
			assert.Equal(t, tc.base.ErrCode, result.ErrCode)
		})
	}
}

func TestErrNo_WithError(t *testing.T) {
	type testCase struct {
		base    ErrNo
		err     error
		wantMsg string
	}

	testCases := map[string]testCase{
		"wrap external error": {
			base:    ServiceErr,
			err:     errors.New("db connection failed"),
			wantMsg: "服务器内部错误: db connection failed",
		},
		"wrap param error": {
			base:    ParamErr,
			err:     errors.New("invalid field"),
			wantMsg: "参数错误: invalid field",
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			result := tc.base.WithError(tc.err)
			assert.Equal(t, tc.wantMsg, result.ErrMsg)
			assert.Equal(t, tc.base.ErrCode, result.ErrCode)
		})
	}
}

func TestConvertErr(t *testing.T) {
	type testCase struct {
		err      error
		wantCode int32
		wantMsg  string
	}

	testCases := map[string]testCase{
		"nil error returns success": {
			err:      nil,
			wantCode: SuccessCode,
			wantMsg:  "成功",
		},
		"errno passthrough": {
			err:      AuthErr,
			wantCode: AuthErrCode,
			wantMsg:  "认证失败",
		},
		"unknown error becomes service error": {
			err:      errors.New("unknown error"),
			wantCode: ServiceErrCode,
			wantMsg:  ServiceErr.ErrMsg,
		},
		"user not exist error passthrough": {
			err:      UserIsNotExistErr,
			wantCode: UserIsNotExistErrCode,
			wantMsg:  "用户不存在",
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			result := ConvertErr(tc.err)
			assert.Equal(t, tc.wantCode, result.ErrCode)
			assert.Equal(t, tc.wantMsg, result.ErrMsg)
		})
	}
}
