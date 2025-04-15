package server

import (
    "fmt"
	"github.com/hofer/nats-mcp/pkg/natsmcp"
    "github.com/mark3labs/mcp-go/server"
	"github.com/nats-io/nats.go"

    log "github.com/sirupsen/logrus"
)

func StartServer(natsUrl string) {
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Create MCP server
	s := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
	)

	natsmcp.RequestTools(nc, s)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}