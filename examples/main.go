package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	log "github.com/sirupsen/logrus"
	"runtime"
	"os"
	"github.com/hofer/nats-mcp/pkg/natstool"
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
	t1 := natstool.NatsMcpTool{
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
	tb := natstool.NewNatsMcpToolBox(t1)

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