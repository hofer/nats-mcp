package natsmcp

import (
	"context"
	"encoding/json"
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
	msg, err := n.nc.Request(n.subject, data, 10*time.Second)
	if err != nil {
		return nil, err
	}
	var result mcp.CallToolResult
	err = json.Unmarshal(msg.Data, &result)
	return &result, err
}

func RequestTools(nc *nats.Conn, s *server.MCPServer) []mcp.Tool {
	tools := []mcp.Tool{}

	doReqAsync(nil, "$SRV.INFO", 0, nc, func(r []byte) {
		var info micro.Info
		json.Unmarshal(r, &info)

		for _, e := range info.Endpoints {
			var tool []mcp.Tool
			json.Unmarshal([]byte(e.Metadata["mcp_tool"]), &tool)
			for _, t := range tool {
				// Add Nats tool handler
				nte := NewNatsToolEndpoint(nc, e.Subject)
				s.AddTool(t, nte.NatsToolHandler)
			}

			tools = append(tools, tool...)
		}
	})
	return tools
}
