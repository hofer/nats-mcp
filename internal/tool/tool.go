package tool

import (
	"context"
	"github.com/hofer/nats-mcp/pkg/natsmcp"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	log "github.com/sirupsen/logrus"
	"time"
)

func StartTool(nc *nats.Conn, command string, env []string, args ...string) (micro.Service, error) {
	ctxClient, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get tool information from the given command (connecting to the tool first):
	mcpClient, err := natsmcp.NewStdioMCPClient(ctxClient, command, env, args...)
	if err != nil {
		return nil, err
	}

	// Get all tools:
	ctxList := context.Background()
	tools, err := mcpClient.ListTools(ctxList, mcp.ListToolsRequest{})
	if err != nil {
		return nil, err
	}

	// Convert the tools found to NatsMcpTools
	var natsMcpTools []natsmcp.NatsMcpTool
	for _, t := range tools.Tools {
		log.Debugf("Tool: %s, Description: %s", t.Name, t.Description)
		mcpTool := natsmcp.NatsMcpTool{
			Tool:    t,
			Handler: mcpClient.CallTool,
		}
		natsMcpTools = append(natsMcpTools, mcpTool)
	}

	// Expose the tools found as a NATS microservice:
	toolBox := natsmcp.NewNatsMcpToolBox(natsMcpTools...)
	srv, err := natsmcp.AddToNewService(nc, toolBox, "NatsService")
	return srv, err
}
