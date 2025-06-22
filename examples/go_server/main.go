package main

import (
	"context"
	"fmt"
	"github.com/hofer/nats-mcp/pkg/natsmcp"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
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
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	return &McpEchoService{
		nc: nc,
	}
}

func (n *McpEchoService) Start() error {
	log.Infof("Starting McpEchoService...")

	// Create tools and corresponding Handlers:
	tools := natsmcp.NatsMcpTool{
		Tool: mcp.NewTool("hello_echo",
			mcp.WithDescription("Say hello to someone"),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Name of the person to greet"),
			)),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			log.Infof("Tool call for hello_echo.")
			name, ok := request.GetArguments()["name"].(string)
			if !ok {
				return mcp.NewToolResultError("name must be a string"), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
		},
	}

	// Add all tools to a toolbox
	toolBox := natsmcp.NewNatsMcpToolBox(tools)

	// Expose all tools in the toolbox as a Nats microservice
	_, err := toolBox.AddToolsAsNatsService(n.nc, "EchoService")
	if err != nil {
		return err
	}

	log.Infof("Service connected at %s", n.nc.ConnectedAddr())
	return nil
}
