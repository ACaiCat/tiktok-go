package ai

import (
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

func TestExistTool(t *testing.T) {
	type testCase struct {
		tools      []string
		queryName  string
		wantExists bool
	}

	testCases := map[string]testCase{
		"registered tool exists": {
			tools:      []string{"tool_a"},
			queryName:  "tool_a",
			wantExists: true,
		},
		"unregistered tool does not exist": {
			tools:      []string{"tool_a"},
			queryName:  "tool_b",
			wantExists: false,
		},
		"empty registry returns false": {
			tools:      []string{},
			queryName:  "tool_a",
			wantExists: false,
		},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			r := toolRegistry{tools: make(map[string]Tool)}
			for _, toolName := range tc.tools {
				r.RegisterTool(newMock(toolName))
			}
			call := openai.ToolCall{Function: openai.FunctionCall{Name: tc.queryName}}
			assert.Equal(t, tc.wantExists, r.ExistTool(call))
		})
	}
}

func TestCallTool(t *testing.T) {
	type testCase struct {
		registeredTools []string
		callName        string
		toolErr         error
		wantErr         bool
		wantContent     string
	}

	testCases := map[string]testCase{
		"call existing tool succeeds": {
			registeredTools: []string{"tool_a"},
			callName:        "tool_a",
			wantErr:         false,
			wantContent:     "result of tool_a",
		},
		"call non-existing tool returns error": {
			registeredTools: []string{"tool_a"},
			callName:        "tool_missing",
			wantErr:         true,
		},
		"tool returns error propagates": {
			registeredTools: []string{"tool_err"},
			callName:        "tool_err",
			toolErr:         assert.AnError,
			wantErr:         true,
		},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			r := toolRegistry{tools: make(map[string]Tool)}
			for _, toolName := range tc.registeredTools {
				m := newMock(toolName)
				m.callErr = tc.toolErr
				r.RegisterTool(m)
			}

			call := openai.ToolCall{Function: openai.FunctionCall{Name: tc.callName}}
			result, err := r.CallTool(call, ToolCallContext{})
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantContent, result.Content)
		})
	}
}

func TestListTools(t *testing.T) {
	type testCase struct {
		tools     []string
		wantCount int
	}

	testCases := map[string]testCase{
		"empty registry lists nothing": {
			tools:     []string{},
			wantCount: 0,
		},
		"single tool listed": {
			tools:     []string{"tool_a"},
			wantCount: 1,
		},
		"multiple tools listed": {
			tools:     []string{"tool_a", "tool_b", "tool_c"},
			wantCount: 3,
		},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			r := toolRegistry{tools: make(map[string]Tool)}
			for _, toolName := range tc.tools {
				r.RegisterTool(newMock(toolName))
			}
			listed := r.ListTools()
			assert.Equal(t, tc.wantCount, len(listed))
			for _, tool := range listed {
				assert.Equal(t, openai.ToolTypeFunction, tool.Type)
				assert.NotEmpty(t, tool.Function.Name)
			}
		})
	}
}
