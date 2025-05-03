package natsmcp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nats-io/nats.go/micro"
	log "github.com/sirupsen/logrus"
)

type NatsMcpTool struct {
	Tool    mcp.Tool
	Handler func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

func NewNatsMcpToolBox(natsTools ...NatsMcpTool) *NatsMcpToolBox {
	return &NatsMcpToolBox{
		natsTools: natsTools,
	}
}

type NatsMcpToolBox struct {
	natsTools []NatsMcpTool
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
	err := json.Unmarshal(request.Data(), &toolRequest)
	if err != nil {
		log.Error(err)
		request.Respond([]byte("Error on tool call (unmarshalling request)"))
		return
	}

	// Find and execute the corresponding tool:
	for _, tt := range t.natsTools {
		if tt.Tool.Name == toolRequest.Params.Name {
			log.Infof("🔧 Calling tool: %s ...", tt.Tool.Name)
			toolResult, tollCallErr := tt.Handler(context.Background(), toolRequest)
			if tollCallErr != nil {
				request.Respond([]byte(fmt.Sprintf("Calling tool %s failed", tt.Tool.Name)))
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

	// We should never reach this point (if the tool name in the request is correct)...
	request.Respond([]byte("Error - maybe the tool could not be found?"))
}

func (t *NatsMcpToolBox) GetSubject() string {
	return "mcp"
}

func (t *NatsMcpToolBox) GetHandlerFunc() micro.Handler {
	return micro.HandlerFunc(t.mcpToolHandler)
}
