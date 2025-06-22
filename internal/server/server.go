package server

import (
	"fmt"
	"github.com/hofer/nats-mcp/pkg/natsmcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"time"
)

func StartServer(nc *nats.Conn, serverName string, serverVersion string, transport string, ssePort string) error {
	// Create MCP server
	s := server.NewMCPServer(serverName, serverVersion)

	// Lookup for all tools exposed via NATS:
	go func() {
		for {
			log.Infof("(Re-)Loading tools from NATS server: %s", nc.ConnectedUrl())
			natsmcp.RequestTools(nc, s)
			time.Sleep(20 * time.Second)
		}
	}()

	// Only check for "sse" since stdio is the default
	if transport == "sse" {
		sseServer := server.NewSSEServer(s, server.WithBaseURL(fmt.Sprintf("http://localhost%s", ssePort)))
		log.Printf("SSE server listening on %s", ssePort)
		return sseServer.Start(ssePort)
	}

	// Start the stdio server
	return server.ServeStdio(s)
}
