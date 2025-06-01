package transport

import (
	"context"
	"encoding/json"
	"fmt"
	mcpTransport "github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nats-io/nats.go/micro"
	log "github.com/sirupsen/logrus"
)

func NewNatsMcpService() *NatsMcpService {
	return &NatsMcpService{}
}

type NatsMcpService struct{}

func (s *NatsMcpService) AddTransportNatsService(srv micro.Service, trans mcpTransport.Interface, serviceName string) error {
	// Add request handler
	rawRoot := srv.AddGroup(RequestGroup)
	err := rawRoot.AddEndpoint(serviceName, micro.HandlerFunc(s.requestFunction(trans)))
	if err != nil {
		return err
	}

	// Add notification handler
	notificationRoot := srv.AddGroup(NotificationGroup)
	err = notificationRoot.AddEndpoint(serviceName, micro.HandlerFunc(s.notificationFunction(trans)))
	return err
}

func (s *NatsMcpService) requestFunction(trans mcpTransport.Interface) func(request micro.Request) {
	return func(request micro.Request) {
		var toolRequest mcpTransport.JSONRPCRequest
		err := json.Unmarshal(request.Data(), &toolRequest)
		if err != nil {
			log.Error(err)
			msg := fmt.Sprintf("Error on mcp_raw request (unmarshalling request): %v", err)
			errorResponse := mcp.NewJSONRPCError(toolRequest.ID, mcp.PARSE_ERROR, msg, request.Data())
			data, _ := json.Marshal(errorResponse)
			request.Respond(data)
			return
		}

		resp, err := trans.SendRequest(context.Background(), toolRequest)
		if err != nil {
			log.Error(err)
			msg := fmt.Sprintf("Error on sending mcp_raw: %v", err)
			errorResponse := mcp.NewJSONRPCError(toolRequest.ID, mcp.INTERNAL_ERROR, msg, "")
			data, _ := json.Marshal(errorResponse)
			request.Respond(data)
			return
		}

		respData, err := json.Marshal(resp)
		request.Respond(respData)
	}
}

func (s *NatsMcpService) notificationFunction(trans mcpTransport.Interface) func(request micro.Request) {
	return func(request micro.Request) {
		var toolRequest mcp.JSONRPCNotification
		err := json.Unmarshal(request.Data(), &toolRequest)
		if err != nil {
			log.Error(err)
			request.Respond([]byte(""))
			return
		}

		err = trans.SendNotification(context.Background(), toolRequest)
		if err != nil {
			log.Error(err)
			request.Respond([]byte(""))
			return
		}
		request.Respond([]byte(""))
	}
}
