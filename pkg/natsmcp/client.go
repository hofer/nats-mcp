package natsmcp

import (
	"context"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewStdioMCPClient(ctx context.Context, command string, env []string, args ...string) (client.MCPClient, error) {
	stdioClient, err := client.NewStdioMCPClient(command, env, args...)
	if err != nil {
		return nil, err
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name: "nats-mcp",
		//Version: "0.0.1",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	_, err = stdioClient.Initialize(ctx, initRequest)
	return stdioClient, err
}
