package server

import (
	"github.com/hofer/nats-mcp/pkg/natsmcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nats-io/nats.go"
)

func StartServer(nc *nats.Conn, serverName string, serverVersion string) error {
	// Create MCP server
	s := server.NewMCPServer(serverName, serverVersion)

	// Lookup for all tools exposed via NATS:
	natsmcp.RequestTools(nc, s)

	// Start the stdio server
	return server.ServeStdio(s)
}
