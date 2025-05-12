package natsmcp

//import (
//	"context"
//	"github.com/mark3labs/mcp-go/client"
//	"github.com/mark3labs/mcp-go/mcp"
//	"github.com/nats-io/nats.go"
//	log "github.com/sirupsen/logrus"
//	"testing"
//)
//
//func TestNewNats(t *testing.T) {
//	nc, err := nats.Connect("nats://100.93.123.116:4222")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	natsTransport := NewNats(nc,"example")
//	natsTransport.Start(context.Background())
//
//	c := client.NewClient(natsTransport)
//	_, err = c.Initialize(context.Background(), mcp.InitializeRequest{})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//
//	toolsResult, err := c.ListTools(context.Background(),mcp.ListToolsRequest{})
//	log.Info(toolsResult)
//
//	//c.CallTool(context.Background(), mcp.CallToolRequest{})
//}