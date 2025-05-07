package natsmcp

import (
	"context"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewNats() *Nats {
	return &Nats{}
}

type Nats struct {

}

func (*Nats) Start(ctx context.Context) error {
	return nil
}

// SendRequest sends a json RPC request and returns the response synchronously.
func (*Nats) 	SendRequest(ctx context.Context, request transport.JSONRPCRequest) (*transport.JSONRPCResponse, error) {
	return nil, nil
}

// SendNotification sends a json RPC Notification to the server.
func (*Nats) SendNotification(ctx context.Context, notification mcp.JSONRPCNotification) error {
	return nil
}

// SetNotificationHandler sets the handler for notifications.
// Any notification before the handler is set will be discarded.
func (*Nats) SetNotificationHandler(handler func(notification mcp.JSONRPCNotification)) {

}

// Close the connection.
func (*Nats) Close() error {
	return nil
}

var _ transport.Interface = (*Nats)(nil)