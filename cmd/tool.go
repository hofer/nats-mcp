package cmd

import (
	"github.com/hofer/nats-mcp/internal/tool"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

var commandStr string

var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Expose tools from a MCP Server via NATS",
	Long: `This command can be used to expose local MCP tools (a MCP server started locally) via NATS. With
just a few simple commands many other MCP servers can be made accessible via NATS.

`,
	Run: func(cmd *cobra.Command, args []string) {
		nc, err := nats.Connect(natsUrl)
		if err != nil {
			log.Fatal(err)
		}
		tool.StartTool(nc, commandStr, args)
		runtime.Goexit()
	},
}

func init() {
	rootCmd.AddCommand(toolCmd)
	toolCmd.Flags().StringVarP(&natsUrl, "url", "u", os.Getenv("NATS_URL"), "URL to the Nats.io server")
	if os.Getenv("NATS_URL") == "" {
		toolCmd.MarkFlagRequired("url")
	}
	toolCmd.Flags().StringVarP(&commandStr, "command", "c", "", "Command to start the local MCP Server")
	toolCmd.MarkFlagRequired("command")
}
