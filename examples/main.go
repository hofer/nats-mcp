package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	log "github.com/sirupsen/logrus"
	"runtime"
	"os"
)

func main() {
	echoService := NewMcpEchoService()
	echoService.Start()
	runtime.Goexit()
}

type McpEchoService struct {
	nc *nats.Conn
}

func NewMcpEchoService() *McpEchoService {
	nc, err := nats.Connect(os.Getenv("NATS_SERVER_URL"))
	if err != nil {
		log.Fatal(err)
	}
	return &McpEchoService{
		nc: nc,
	}
}

func (n *McpEchoService) Start() {
	log.Infof("Starting McpEchoService...")

	srv, err := micro.AddService(n.nc, micro.Config{
		Name:        "McpEchoService",
		Version:     "0.0.2",
		Description: "Simple MCP echo service to test MCP via Nats.",
	})
	if err != nil {
		log.Fatal(err)
	}
	//defer srv.Stop()

	root := srv.AddGroup("echo")

	// Create tools and corresponding Handlers:
	t1 := NatsMcpTool{
		Tool: mcp.NewTool("hello_echo",
			mcp.WithDescription("Say hello to someone"),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Name of the person to greet"),
			)),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			log.Infof("Tool call for hello_echo.")
			name, ok := request.Params.Arguments["name"].(string)
			if !ok {
				return mcp.NewToolResultError("name must be a string"), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
		},
	}

	// Add tools to toolbox
	tb := NewNatsMcpToolBox(t1)

	// Mcp Echo Endpoint
	err = root.AddEndpoint(
		tb.GetSubject(),
		tb.GetHandlerFunc(),
		micro.WithEndpointMetadata(map[string]string{
			"mcp_tool": tb.CreateMcpToolMetadata(),
		}))
	if err != nil {
		log.Fatal(err)
	}

	n.nc.ConnectedAddr()
}

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
