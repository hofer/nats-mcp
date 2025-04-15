package natsmcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go/micro"
)

type NatsMcpTool struct {
	Tool    mcp.Tool
	Handler func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

type NatsMcpToolBox struct {
	natsTools []NatsMcpTool
}

func NewNatsMcpToolBox(natsTools ...NatsMcpTool) *NatsMcpToolBox {
	return &NatsMcpToolBox{
		natsTools: natsTools,
	}
}

func (t *NatsMcpToolBox) CreateMcpToolMetadata() string {
	tools := []mcp.Tool{}
	for _, t := range t.natsTools {
		tools = append(tools, t.Tool)
	}
	jsonStr, _ := json.Marshal(tools)
	return string(jsonStr)
}

func (t *NatsMcpToolBox) mcpToolHandler(request micro.Request) {
	// Get the tool request:
	var toolRequest mcp.CallToolRequest
	json.Unmarshal(request.Data(), &toolRequest)

	for _, tt := range t.natsTools {
		if tt.Tool.Name == toolRequest.Params.Name {
			toolResult, err := tt.Handler(context.Background(), toolRequest)
			if err != nil {
				// TODO: Implement error
				request.Respond([]byte("hello"))
				return
			} else {
				toolResultJson, err := json.Marshal(toolResult)
				if err != nil {
					request.Respond([]byte("Error"))
					return
				}
				request.Respond(toolResultJson)
				return
			}
		}
	}

	request.Respond([]byte("Error"))
}

func (t *NatsMcpToolBox) GetSubject() string {
	return "mcp"
}

func (t *NatsMcpToolBox) GetHandlerFunc() micro.Handler {
	return micro.HandlerFunc(t.mcpToolHandler)
}
