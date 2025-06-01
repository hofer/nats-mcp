package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewNats(t *testing.T) {
	// Create an embedded Nats server for testing...
	nc := StartEmbeddedNats()

	// Create a simple mock response
	_, err := nc.Subscribe(fmt.Sprintf("%s.%s", RequestGroup, "nats-mcp-transport-test"), func(msg *nats.Msg) {
		data, err := json.Marshal(transport.JSONRPCResponse{
			Result: []byte("{ \"name\": \"TestResponse\"}"),
		})
		assert.NoError(t, err)
		msg.Respond(data)
	})
	if err != nil {
		t.Errorf("Error on subscribe to nats server: %v", err)
	}

	natsTransport := NewNats(nc, "nats-mcp-transport-test")
	natsTransport.Start(context.Background())

	resp, err := natsTransport.SendRequest(context.Background(), transport.JSONRPCRequest{
		Method: "mcp_raw.nats-mcp-transport-test",
	})
	if err != nil {
		t.Errorf("Error on sending request: %v", err)
	}

	assert.Contains(t, string(resp.Result), "TestResponse")
}

func StartEmbeddedNats() *nats.Conn {
	opts := &server.Options{}

	// Initialize new server with options
	ns, err := server.NewServer(opts)
	if err != nil {
		panic(err)
	}

	// Start the server via goroutine
	go ns.Start()

	// Wait for server to be ready for connections
	if !ns.ReadyForConnections(4 * time.Second) {
		panic("not ready for connection")
	}

	// Connect to server
	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		panic(err)
	}
	return nc
}
