package client

import (
	"github.com/hofer/nats-mcp/pkg/natsmcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

func ListTools(nc *nats.Conn) error {
	// Create MCP server
	s := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
	)

	tools := natsmcp.RequestTools(nc, s)

	// Add implementation to get tools:
	for _, t := range tools {
		log.Infof("🛠️ %s: %s", t.Name, t.Description)
	}

	return nil
}
