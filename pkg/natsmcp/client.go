package natsmcp

import (
	"context"
	"github.com/hofer/nats-mcp/pkg/natsmcp/client"
	mcpClient "github.com/mark3labs/mcp-go/client"
	mcpTransport "github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nats-io/nats.go"
)

func NewStdioMCPClient(ctx context.Context, command string, env []string, args ...string) (mcpClient.MCPClient, mcpTransport.Interface, error) {
	stdioClient, err := mcpClient.NewStdioMCPClient(command, env, args...)
	if err != nil {
		return nil, nil, err
	}

	err = stdioClient.GetTransport().Start(context.Background())
	if err != nil {
		return nil, nil, err
	}

	initRequest := createInitRequest()
	_, err = stdioClient.Initialize(ctx, initRequest)
	return stdioClient, stdioClient.GetTransport(), err
}

func NewSSEMCPClient(ctx context.Context, baseUrl string, options ...mcpTransport.ClientOption) (mcpClient.MCPClient, mcpTransport.Interface, error) {
	sseClient, err := mcpClient.NewSSEMCPClient(baseUrl, options...)
	if err != nil {
		return nil, nil, err
	}

	err = sseClient.GetTransport().Start(context.Background())
	if err != nil {
		return nil, nil, err
	}

	initRequest := createInitRequest()
	_, err = sseClient.Initialize(ctx, initRequest)
	return sseClient, sseClient.GetTransport(), err
}

func NewNatsMCPClient(nc *nats.Conn, ctx context.Context, subject string) (mcpClient.MCPClient, mcpTransport.Interface, error) {
	natsClient, err := client.NewNatsMCPClient(nc, subject)
	if err != nil {
		return nil, nil, err
	}

	err = natsClient.GetTransport().Start(context.Background())
	if err != nil {
		return nil, nil, err
	}

	initRequest := createInitRequest()
	_, err = natsClient.Initialize(ctx, initRequest)
	return natsClient, natsClient.GetTransport(), err
}

func createInitRequest() mcp.InitializeRequest {
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "nats-mcp",
		Version: "0.0.1",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}
	return initRequest
}
