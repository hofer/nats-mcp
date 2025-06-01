package natsmcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/klauspost/compress/s2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
	"sync"
	"time"
)

type NatsMcpTool struct {
	Tool    mcp.Tool
	Handler func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

func NewNatsMcpToolBox(natsTools ...NatsMcpTool) *natsMcpToolBox {
	return &natsMcpToolBox{
		natsTools: natsTools,
	}
}

type natsMcpToolBox struct {
	natsTools []NatsMcpTool
}

func (t *natsMcpToolBox) CreateMcpToolMetadata() string {
	tools := []mcp.Tool{}
	for _, t := range t.natsTools {
		tools = append(tools, t.Tool)
	}
	jsonStr, _ := json.Marshal(tools)
	return string(jsonStr)
}

func (t *natsMcpToolBox) mcpToolHandler(request micro.Request) {
	// Get the tool request:
	var toolRequest mcp.CallToolRequest
	err := json.Unmarshal(request.Data(), &toolRequest)
	if err != nil {
		log.Error(err)
		request.Respond([]byte("Error on tool call (unmarshalling request)"))
		return
	}

	// Find and execute the corresponding tool:
	for _, tt := range t.natsTools {
		if tt.Tool.Name == toolRequest.Params.Name {
			log.Infof("🔧 Calling tool: %s ...", tt.Tool.Name)
			toolResult, tollCallErr := tt.Handler(context.Background(), toolRequest)
			if tollCallErr != nil {
				log.Error(tollCallErr)
				request.Respond([]byte(fmt.Sprintf("Calling tool %s failed", tt.Tool.Name)))
				return
			} else {
				toolResultJson, err := json.Marshal(toolResult)
				if err != nil {
					log.Error(err)
					request.Respond([]byte("Error"))
					return
				}
				log.Debug(string(toolResultJson))
				request.Respond(toolResultJson)
				return
			}
		}
	}

	// We should never reach this point (if the tool name in the request is correct)...
	request.Respond([]byte("Error - maybe the tool could not be found?"))
}

func (t *natsMcpToolBox) GetSubject() string {
	return "mcp"
}

func (t *natsMcpToolBox) GetHandlerFunc() micro.Handler {
	return micro.HandlerFunc(t.mcpToolHandler)
}

func (t *natsMcpToolBox) AddToolsAsNatsService(nc *nats.Conn, serviceName string) (micro.Service, error) {
	srv, err := micro.AddService(nc, micro.Config{
		Name:        fmt.Sprintf("%sMCP", serviceName),
		Version:     "0.0.2",
		Description: fmt.Sprintf("MCP service for %s exposing tools via NATS microservices.", serviceName),
	})
	if err != nil {
		return srv, err
	}
	//defer srv.Stop()

	// MCP Tool Endpoint
	toolRoot := srv.AddGroup("mcp_tool")
	err = toolRoot.AddEndpoint(
		serviceName,
		t.GetHandlerFunc(),
		micro.WithEndpointMetadata(map[string]string{
			"mcp_tool": t.CreateMcpToolMetadata(),
		}))
	return srv, err
}

// doReqAsync serializes and sends a request to the given subject and handles multiple responses.
// The value of the `waitFor` may shorten the interval during which responses are gathered:
//
//	waitFor < 0  : listen for responses for the full timeout interval
//	waitFor == 0 : (adaptive timeout), after each response, wait a short amount of time for more, then stop
//	waitFor > 0  : stops listening before the timeout if the given number of responses are received
func doReqAsync(req any, subj string, waitFor int, nc *nats.Conn, cb func([]byte)) error {
	ctx := context.Background()
	timeout := 10 * time.Second

	jreq := []byte("{}")
	var err error

	if req != nil {
		switch val := req.(type) {
		case string:
			jreq = []byte(val)
		default:
			jreq, err = json.Marshal(req)
			if err != nil {
				return err
			}
		}
	}

	var (
		mu       sync.Mutex
		ctr      = 0
		finisher *time.Timer
	)

	// Set deadline, max amount of time this function waits for responses
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Activate "adaptive timeout". Finisher may trigger early termination
	if waitFor == 0 {
		// First response can take up to Timeout to arrive
		finisher = time.NewTimer(timeout)
		go func() {
			select {
			case <-finisher.C:
				cancel()
			case <-ctx.Done():
				return
			}
		}()
	}

	errs := make(chan error)
	sub, err := nc.Subscribe(nc.NewRespInbox(), func(m *nats.Msg) {
		mu.Lock()
		defer mu.Unlock()

		data := m.Data
		//compressed := false
		if m.Header.Get("Content-Encoding") == "snappy" {
			//compressed = true
			ud, err := io.ReadAll(s2.NewReader(bytes.NewBuffer(data)))
			if err != nil {
				errs <- err
				return
			}
			data = ud
		}

		// If adaptive timeout is active, set deadline for next response
		if finisher != nil {
			// Stop listening and return if no further responses arrive within this interval
			finisher.Reset(300 * time.Millisecond)
		}

		if m.Header.Get("Status") == "503" {
			errs <- nats.ErrNoResponders
			return
		}

		cb(data)
		ctr++

		// Stop listening if the requested number of responses have been received
		if waitFor > 0 && ctr == waitFor {
			cancel()
		}
	})
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	if waitFor > 0 {
		sub.AutoUnsubscribe(waitFor)
	}

	msg := nats.NewMsg(subj)
	msg.Data = jreq
	if subj != "$SYS.REQ.SERVER.PING" && !strings.HasPrefix(subj, "$SYS.REQ.ACCOUNT") {
		msg.Header.Set("Accept-Encoding", "snappy")
	}
	msg.Reply = sub.Subject

	err = nc.PublishMsg(msg)
	if err != nil {
		return err
	}

	select {
	case err = <-errs:
		if err == nats.ErrNoResponders && strings.HasPrefix(subj, "$SYS") {
			return fmt.Errorf("server request failed, ensure the account used has system privileges and appropriate permissions")
		}

		return err
	case <-ctx.Done():
	}

	return nil
}
