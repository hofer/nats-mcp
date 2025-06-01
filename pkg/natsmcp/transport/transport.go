package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nats-io/nats.go"
	"time"
)

const RequestGroup = "mcp_raw"
const NotificationGroup = "mcp_notification"

var _ transport.Interface = (*Nats)(nil)

func NewNats(nc *nats.Conn, subject string) *Nats {
	return &Nats{
		nc:      nc,
		subject: subject,
	}
}

type Nats struct {
	nc      *nats.Conn
	subject string
}

func (t *Nats) Start(ctx context.Context) error {
	return nil
}

// SendRequest sends a json RPC request and returns the response synchronously.
func (t *Nats) SendRequest(ctx context.Context, request transport.JSONRPCRequest) (*transport.JSONRPCResponse, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	remainingDuration := time.Second * 30
	deadline, ok := ctx.Deadline()
	if ok {
		remainingDuration = time.Until(deadline)
	}

	msg, err := t.nc.Request(fmt.Sprintf("%s.%s", RequestGroup, t.subject), data, remainingDuration)
	if msg == nil {
		return nil, fmt.Errorf("no response from server")
	}

	var result transport.JSONRPCResponse
	err = json.Unmarshal(msg.Data, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSONRPCResponse: %v, data: %s", err, msg.Data)
	}

	return &result, err
}

// SendNotification sends a json RPC Notification to the server.
func (t *Nats) SendNotification(ctx context.Context, notification mcp.JSONRPCNotification) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	remainingDuration := time.Second * 30
	deadline, ok := ctx.Deadline()
	if ok {
		remainingDuration = time.Until(deadline)
	}

	msg, err := t.nc.Request(fmt.Sprintf("%s.%s", NotificationGroup, t.subject), data, remainingDuration)

	if len(msg.Data) > 0 {
		return fmt.Errorf("Error on notification %s", string(msg.Data))
	}

	return err
}

// SetNotificationHandler sets the handler for notifications.
// Any notification before the handler is set will be discarded.
func (t *Nats) SetNotificationHandler(handler func(notification mcp.JSONRPCNotification)) {

}

// Close the connection.
func (t *Nats) Close() error {
	return nil
}
