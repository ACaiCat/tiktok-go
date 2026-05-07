package ai

import (
	"context"

	service "github.com/ACaiCat/tiktok-go/biz/service/user"
)

type JwchLoginInput struct {
	UserID int64 `json:"user_id" jsonschema:"required,description=聊天用户ID"`
}

type JwchLoginOutput struct {
	JwchID     string `json:"jwch_id" jsonschema:"required,description=教务处ID"`
	JwchCookie string `json:"jwch_cookie" jsonschema:"required,description=教务处会话Cookie"`
}

func jwchLoginTool() LocalTool[JwchLoginInput, JwchLoginOutput] {
	return LocalTool[JwchLoginInput, JwchLoginOutput]{
		Name:        "internal_jwch_login",
		Description: "Logs in to the Jwch system using a user ID. Use this as the required login tool when cookies or credentials are unavailable in the context.",
		Func:        JwchLogin,
	}
}

func JwchLogin(input JwchLoginInput) (JwchLoginOutput, error) {
	id, cookie, err := service.NewUserService(context.Background()).GetJwchIdentifierAndCookies(input.UserID)
	if err != nil {
		return JwchLoginOutput{}, err
	}

	return JwchLoginOutput{
		JwchID:     id,
		JwchCookie: cookie,
	}, nil
}
