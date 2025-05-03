package natsmcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/klauspost/compress/s2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"io"
	"strings"
	"sync"
	"time"
)

func AddToNewService(nc *nats.Conn, toolBox *NatsMcpToolBox, serviceName string) (micro.Service, error) {
	srv, err := micro.AddService(nc, micro.Config{
		Name:        fmt.Sprintf("%sMCP", serviceName),
		Version:     "0.0.2",
		Description: "MCP service exposing MCP tools via Nats.",
	})
	if err != nil {
		return srv, err
	}
	//defer srv.Stop()

	root := srv.AddGroup("mcp")

	// Mcp Echo Endpoint
	err = root.AddEndpoint(
		toolBox.GetSubject(),
		toolBox.GetHandlerFunc(),
		micro.WithEndpointMetadata(map[string]string{
			"mcp_tool": toolBox.CreateMcpToolMetadata(),
		}))

	return srv, nil
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
