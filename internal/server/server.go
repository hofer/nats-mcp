package server

import (
	"fmt"
	"github.com/hofer/nats-mcp/pkg/natsmcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

func StartServer(nc *nats.Conn, serverName string, serverVersion string, transport string, ssePort string) error {
	// Create MCP server
	s := server.NewMCPServer(serverName, serverVersion)

	// Lookup for all tools exposed via NATS:
	natsmcp.RequestTools(nc, s)

	// Only check for "sse" since stdio is the default
	if transport == "sse" {
		sseServer := server.NewSSEServer(s, server.WithBaseURL(fmt.Sprintf("http://localhost%s", ssePort)))
		log.Printf("SSE server listening on %s", ssePort)
		return sseServer.Start(ssePort)
	}

	// Start the stdio server
	return server.ServeStdio(s)
}
