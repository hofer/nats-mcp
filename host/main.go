package main

import (
	"context"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"os/signal"
	"time"
	log "github.com/sirupsen/logrus"
)

func main() {
	mcpClient, err := createClient()
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


	shutdownhandler()
}

func createClient() (*client.StdioMCPClient, error) {
	client, err := client.NewStdioMCPClient("../server/mcp-server", []string{})
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

func shutdownhandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for range c {
		log.Info("about to shutdown the app.....")
		os.Exit(0)
		// sig is a ^C, handle it
	}
}
