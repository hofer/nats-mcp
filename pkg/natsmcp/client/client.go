package client

import (
	"context"
	"fmt"
	"github.com/hofer/nats-mcp/pkg/natsmcp/transport"
	"github.com/mark3labs/mcp-go/client"
	"github.com/nats-io/nats.go"
)

func NewNatsMCPClient(nc *nats.Conn, subject string) (*client.Client, error) {
	natsTransport := transport.NewNats(nc, subject)
	err := natsTransport.Start(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to start stdio transport: %w", err)
	}
	return client.NewClient(natsTransport), nil
}
