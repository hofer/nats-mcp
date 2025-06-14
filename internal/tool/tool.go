package tool

import (
	"context"
	"github.com/hofer/nats-mcp/pkg/natsmcp"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	log "github.com/sirupsen/logrus"
	"time"
)

func StartStdioTools(nc *nats.Conn, serverName string, command string, env []string, args ...string) (micro.Service, error) {
	ctxClient, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get tool information from the given command (connecting to the tool first):
	mcpClient, _, err := natsmcp.NewStdioMCPClient(ctxClient, command, env, args...)
	if err != nil {
		return nil, err
	}

	return startTools(nc, mcpClient, serverName)
}

func StartSseTools(nc *nats.Conn, serverName string, baseUrl string) (micro.Service, error) {
	ctxClient, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get tool information from the given command (connecting to the tool first):
	mcpClient, _, err := natsmcp.NewSSEMCPClient(ctxClient, baseUrl)
	if err != nil {
		return nil, err
	}

	return startTools(nc, mcpClient, serverName)
}

func startTools(nc *nats.Conn, mcpClient client.MCPClient, serverName string) (micro.Service, error) {
	// Get all tools:
	ctxList := context.Background()
	tools, err := mcpClient.ListTools(ctxList, mcp.ListToolsRequest{})
	if err != nil {
		return nil, err
	}

	// Convert the tools found to NatsMcpTools
	var natsMcpTools []natsmcp.NatsMcpTool
	for _, t := range tools.Tools {
		log.Debugf("Server: %s, Tool: %s, Description: %s", serverName, t.Name, t.Description)
		mcpTool := natsmcp.NatsMcpTool{
			Tool:    t,
			Handler: mcpClient.CallTool,
		}
		natsMcpTools = append(natsMcpTools, mcpTool)
	}

	// Expose the tools found as a NATS microservice:
	log.Infof("🚀 Starting MCP Server and exposing tools...")
	toolBox := natsmcp.NewNatsMcpToolBox(natsMcpTools...)
	srv, err := toolBox.AddToolsAsNatsService(nc, serverName)
	if err != nil {
		return srv, err
	}

	// TODO: Do we want to expose tools as Raw services as well?
	// Adding Nats Service for direct tool calling
	// err = toolBox.AddTransportNatsService(srv, transport, serverName)
	return srv, err
}
