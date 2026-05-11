package ai

import (
	"errors"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

type addInput struct {
	A int `json:"a"`
	B int `json:"b"`
}

type addOutput struct {
	Sum int `json:"sum"`
}

func newAddTool(authorize func(ToolCallContext, addInput) error, fn func(addInput) (addOutput, error)) LocalTool[addInput, addOutput] {
	return LocalTool[addInput, addOutput]{
		Name:        "add",
		Description: "adds two numbers",
		Authorize:   authorize,
		Func:        fn,
	}
}

func TestLocalToolGetName(t *testing.T) {
	type testCase struct {
		name string
	}

	testCases := map[string]testCase{
		"returns correct name": {name: "my_tool"},
		"returns empty name":   {name: ""},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			tool := LocalTool[addInput, addOutput]{Name: tc.name}
			assert.Equal(t, tc.name, tool.GetName())
		})
	}
}

func TestLocalToolGetDescription(t *testing.T) {
	type testCase struct {
		description string
	}

	testCases := map[string]testCase{
		"returns correct description": {description: "does something"},
		"returns empty description":   {description: ""},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			tool := LocalTool[addInput, addOutput]{Description: tc.description}
			assert.Equal(t, tc.description, tool.GetDescription())
		})
	}
}

func TestLocalToolGetParametersDef(t *testing.T) {
	PatchConvey("reflects schema from input type", t, func() {
		tool := LocalTool[addInput, addOutput]{}
		params := tool.GetParametersDef()
		assert.NotNil(t, params)
	})
}

func TestLocalToolCallTool(t *testing.T) {
	type testCase struct {
		args      string
		authorize func(ToolCallContext, addInput) error
		fn        func(addInput) (addOutput, error)
		wantErr   bool
		checkMsg  func(t *testing.T, msg *openai.ChatCompletionMessage)
	}

	testCases := map[string]testCase{
		"success returns json result": {
			args:      `{"a":1,"b":2}`,
			authorize: nil,
			fn: func(in addInput) (addOutput, error) {
				return addOutput{Sum: in.A + in.B}, nil
			},
			wantErr: false,
			checkMsg: func(t *testing.T, msg *openai.ChatCompletionMessage) {
				assert.Equal(t, openai.ChatMessageRoleTool, msg.Role)
				assert.Contains(t, msg.Content, "sum")
			},
		},
		"invalid json arguments returns error": {
			args:    `not json`,
			fn:      func(in addInput) (addOutput, error) { return addOutput{}, nil },
			wantErr: true,
		},
		"authorize non-errno error wraps as AuthErr": {
			args: `{"a":1,"b":2}`,
			authorize: func(_ ToolCallContext, _ addInput) error {
				return errors.New("forbidden")
			},
			fn:      func(in addInput) (addOutput, error) { return addOutput{}, nil },
			wantErr: true,
			checkMsg: func(t *testing.T, msg *openai.ChatCompletionMessage) {
				// msg is nil when authorize fails
			},
		},
		"authorize errno error passes through": {
			args: `{"a":1,"b":2}`,
			authorize: func(_ ToolCallContext, _ addInput) error {
				return errno.AuthErr
			},
			fn:      func(in addInput) (addOutput, error) { return addOutput{}, nil },
			wantErr: true,
		},
		"func error returns isError result not an error": {
			args:      `{"a":1,"b":2}`,
			authorize: nil,
			fn: func(in addInput) (addOutput, error) {
				return addOutput{}, errors.New("compute failed")
			},
			wantErr: false,
			checkMsg: func(t *testing.T, msg *openai.ChatCompletionMessage) {
				assert.Contains(t, msg.Content, "call tool error")
			},
		},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			tool := newAddTool(tc.authorize, tc.fn)
			call := openai.ToolCall{
				ID:       "call_1",
				Function: openai.FunctionCall{Name: "add", Arguments: tc.args},
			}

			msg, err := tool.CallTool(call, ToolCallContext{})
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, msg)
			if tc.checkMsg != nil {
				tc.checkMsg(t, msg)
			}
		})
	}
}
