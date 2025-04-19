package client

import (
	"github.com/hofer/nats-mcp/pkg/natsmcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

func ListTools(natsUrl string) error {
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return err
	}

	// Create MCP server
	s := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
	)

	tools := natsmcp.RequestTools(nc, s)

	// Add implementation to get tools:
	for _, t := range tools {
		log.Info("%s: %s", t.Name, t.Description)
	}

	return nil
}
