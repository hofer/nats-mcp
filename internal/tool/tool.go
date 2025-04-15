package tool

import (
	"context"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"time"
	log "github.com/sirupsen/logrus"
)

func StartTool(natsUrl string, command string, args []string) {
	mcpClient, err := createToolClient(command)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	tools, err := mcpClient.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tools.Tools {
		log.Infof("Tool: %s, Description: %s", t.Name, t.Description)
	}
}

func createToolClient(command string) (*client.StdioMCPClient, error) {
	client, err := client.NewStdioMCPClient(command, []string{})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("Initializing server...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "natsmicromcphost",
		Version: "0.0.1",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	_, err = client.Initialize(ctx, initRequest)
	return client, err
}