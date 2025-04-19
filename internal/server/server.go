package server

import (
	"github.com/hofer/nats-mcp/pkg/natsmcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nats-io/nats.go"
)

func StartServer(natsUrl string) error {
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return err
	}

	// Create MCP server
	s := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
	)

	natsmcp.RequestTools(nc, s)

	// Start the stdio server
	err = server.ServeStdio(s)
	return err
}
