package natsmcp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"time"
)

func NewNatsToolEndpoint(nc *nats.Conn, subject string) *NatsToolEndpoint {
	return &NatsToolEndpoint{
		nc:      nc,
		subject: subject,
	}
}

type NatsToolEndpoint struct {
	nc      *nats.Conn
	subject string
}

func (n *NatsToolEndpoint) NatsToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	remainingDuration := time.Second * 30
	deadline, ok := ctx.Deadline()
	if ok {
		remainingDuration = time.Until(deadline)
	}

	msg, err := n.nc.Request(n.subject, data, remainingDuration)
	if err != nil {
		return nil, err
	}

	rawMessage := (*json.RawMessage)(&msg.Data)
	result, err := mcp.ParseCallToolResult(rawMessage)
	if err != nil {
		return nil, fmt.Errorf("error parsing tool result: %v, data: %s", err, msg.Data)
	}

	return result, err
}

func RequestTools(nc *nats.Conn, mcpServer *server.MCPServer) ([]mcp.Tool, error) {
	tools := []mcp.Tool{}
	serverTools := []server.ServerTool{}
	err := doReqAsync(nil, "$SRV.INFO", 0, nc, func(r []byte) {
		var info micro.Info
		json.Unmarshal(r, &info)

		for _, e := range info.Endpoints {
			var tool []mcp.Tool
			json.Unmarshal([]byte(e.Metadata["mcp_tool"]), &tool)
			for _, t := range tool {
				// Add Nats tool handler
				nte := NewNatsToolEndpoint(nc, e.Subject)
				serverTools = append(serverTools, server.ServerTool{
					Tool:    t,
					Handler: nte.NatsToolHandler,
				})
			}

			tools = append(tools, tool...)
		}
	})

	if err != nil {
		return []mcp.Tool{}, err
	}

	mcpServer.SetTools(serverTools...)
	return tools, nil
}
