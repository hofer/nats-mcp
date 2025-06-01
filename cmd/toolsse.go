package cmd

import (
	"github.com/hofer/nats-mcp/internal/tool"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

var sseUrl string

var toolSseCmd = &cobra.Command{
	Use:   "sse",
	Short: "Expose tools from a local MCP Server (SSE) via NATS",
	Long: `This command can be used to expose local MCP tools (a MCP server started locally) via NATS. With
just a few simple commands many different MCP servers can be made accessible via NATS.
`,
	Run: func(cmd *cobra.Command, args []string) {
		nc, err := nats.Connect(natsUrl)
		if err != nil {
			log.Fatal(err)
		}

		err = StartSseTool(nc, toolServerName, sseUrl)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Waiting for incoming tool calls...")
		runtime.Goexit()
	},
}

func init() {
	toolCmd.AddCommand(toolSseCmd)
	toolSseCmd.Flags().StringVarP(&natsUrl, "url", "u", os.Getenv("NATS_URL"), "URL to the Nats.io server")
	if os.Getenv("NATS_URL") == "" {
		toolCmd.MarkFlagRequired("url")
	}

	toolSseCmd.Flags().StringVarP(&toolServerName, "name", "n", "", "Server name (if used with commandline args)")
	toolSseCmd.MarkFlagRequired("name")

	toolSseCmd.Flags().StringVarP(&sseUrl, "sseUrl", "s", "", "Url to the MCP SSE Server")
	toolSseCmd.MarkFlagRequired("sseUrl")
}

func StartSseTool(nc *nats.Conn, serverName string, baseUrl string) error {
	_, err := tool.StartSseTools(nc, serverName, baseUrl)
	return err
}
