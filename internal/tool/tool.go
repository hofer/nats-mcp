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

func StartTool(nc *nats.Conn, command string, env []string, args ...string) (micro.Service, error) {
	// Get tool information from the given command (connecting to the tool first):
	mcpClient, err := createToolClient(command, env, args...)
	if err != nil {
		return nil, err
	}

	// Get all tools:
	ctx := context.Background()
	tools, err := mcpClient.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return nil, err
	}

	// Convert the tools found to NatsMcpTools
	var natsMcpTools []natsmcp.NatsMcpTool
	for _, t := range tools.Tools {
		log.Debugf("Tool: %s, Description: %s", t.Name, t.Description)
		mcpTool := natsmcp.NatsMcpTool{
			Tool: t,
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return nil, nil
			},
		}
		natsMcpTools = append(natsMcpTools, mcpTool)
	}

	// Expose the tools found as a NATS microservice:
	toolBox := natsmcp.NewNatsMcpToolBox(natsMcpTools...)
	srv, err := natsmcp.AddToNewService(nc, toolBox)
	return srv, err
}

func createToolClient(command string, env []string, args ...string) (*client.StdioMCPClient, error) {
	stdioClient, err := client.NewStdioMCPClient(command, env, args...)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("Initializing server...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "natsmicromcphost",
		Version: "0.0.1",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	_, err = stdioClient.Initialize(ctx, initRequest)
	return stdioClient, err
}
