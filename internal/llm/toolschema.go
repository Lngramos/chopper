package llm

import "fmt"

type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

func (tc ToolCall) Debug() {
	fmt.Printf("[DEBUG] ToolCall detected: name=%s, args=%+v\n", tc.Name, tc.Arguments)
}

type ToolCheck struct {
	ToolCall  *ToolCall  `json:"tool_call"`
	ToolCalls []ToolCall `json:"tool_calls"`
}

func (tc ToolCheck) Debug() {
	if tc.ToolCall != nil {
		fmt.Println("[DEBUG] Single tool_call present")
		tc.ToolCall.Debug()
	} else if len(tc.ToolCalls) > 0 {
		fmt.Printf("[DEBUG] Multiple tool_calls present: count=%d\n", len(tc.ToolCalls))
		for i, call := range tc.ToolCalls {
			fmt.Printf("[DEBUG] tool_calls[%d]:\n", i)
			call.Debug()
		}
	} else {
		fmt.Println("[DEBUG] No tool calls present")
	}
}
