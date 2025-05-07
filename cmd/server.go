package cmd

import (
	"github.com/hofer/nats-mcp/internal/server"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var serverName string
var serverVersion string
var natsUrl string
var transport string
var ssePort string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start an MCP Server exposing all NATS MCP tools",
	Long: `This command will start a MCP Server and make available all tools which are defined
as NATS microservices as tools.

The server is accessible via stdio only. 
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Connecting to the Nats.io Server: %s", natsUrl)
		nc, err := nats.Connect(natsUrl)
		if err != nil {
			log.Fatal(err)
		}
		err = server.StartServer(nc, serverName, serverVersion, transport, ssePort)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&serverName, "name", "n", "NatsMcp Server 🚀", "MCP server name")
	serverCmd.Flags().StringVarP(&serverVersion, "version", "v", "1.0.0", "MCP server version")

	serverCmd.Flags().StringVarP(&natsUrl, "url", "u", os.Getenv("NATS_URL"), "URL to the Nats.io server")
	if os.Getenv("NATS_URL") == "" {
		serverCmd.MarkFlagRequired("url")
	}

	serverCmd.Flags().StringVarP(&transport, "transport", "t", "stdio", "Transport can either be 'sse' or 'stdio'")
	serverCmd.Flags().StringVarP(&ssePort, "port", "p", ":8080", "Port the SSE server should listen on")
}
