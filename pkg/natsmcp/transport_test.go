package natsmcp

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestNewNats(t *testing.T) {
	nc := StartEmbeddedNats()

	natsTransport := NewNats(nc, "example")
	natsTransport.Start(context.Background())

	StartTool(nc, natsTransport)

	client.NewClient(natsTransport)

	// TODO:
	//c := client.NewClient(natsTransport)
	//_, err := c.Initialize(context.Background(), mcp.InitializeRequest{})
	//if err != nil {
	//	t.Error(err)
	//}

	//toolsResult, err := c.ListTools(context.Background(), mcp.ListToolsRequest{})
	//log.Info(toolsResult)

	//c.CallTool(context.Background(), mcp.CallToolRequest{})
}

func StartEmbeddedNats() *nats.Conn {
	opts := &server.Options{}

	// Initialize new server with options
	ns, err := server.NewServer(opts)
	if err != nil {
		panic(err)
	}

	// Start the server via goroutine
	go ns.Start()

	// Wait for server to be ready for connections
	if !ns.ReadyForConnections(4 * time.Second) {
		panic("not ready for connection")
	}

	// Connect to server
	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		panic(err)
	}
	return nc
}

func StartTool(nc *nats.Conn, transport *Nats) error {
	log.Infof("Starting McpEchoService...")

	// Create tools and corresponding Handlers:
	tools := NatsMcpTool{
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

	serviceName := "EchoService"

	// Add all tools to a toolbox
	toolBox := NewNatsMcpToolBox(tools)

	// Expose all tools in the toolbox as a Nats microservice
	srv, err := toolBox.AddToolsAsNatsService(nc, serviceName)
	if err != nil {
		return err
	}

	toolBox.AddTransportNatsService(srv, transport, serviceName)

	log.Infof("Service connected at %s", nc.ConnectedAddr())
	return nil
}
