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
var toolServerName string
var environment []string

var toolStdioCmd = &cobra.Command{
	Use:   "stdio",
	Short: "Expose tools from a local MCP Server (Stdio) via NATS",
	Long: `This command can be used to expose local MCP tools (a MCP server started locally) via NATS. With
just a few simple commands many different MCP servers can be made accessible via NATS.
`,
	Run: func(cmd *cobra.Command, args []string) {
		nc, err := nats.Connect(natsUrl)
		if err != nil {
			log.Fatal(err)
		}

		err = StartStdioTool(nc, toolServerName, commandStr, environment, args...)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Waiting for incoming tool calls...")
		runtime.Goexit()
	},
}

func init() {
	toolCmd.AddCommand(toolStdioCmd)
	toolStdioCmd.Flags().StringVarP(&natsUrl, "url", "u", os.Getenv("NATS_URL"), "URL to the Nats.io server")
	if os.Getenv("NATS_URL") == "" {
		toolCmd.MarkFlagRequired("url")
	}

	toolStdioCmd.Flags().StringVarP(&toolServerName, "name", "n", "", "Server name (if used with commandline args)")
	toolStdioCmd.MarkFlagRequired("name")

	toolStdioCmd.Flags().StringVarP(&commandStr, "command", "c", "", "Command to start the local MCP Server")
	toolStdioCmd.MarkFlagRequired("command")

	toolStdioCmd.Flags().StringArrayVarP(&environment, "env", "e", []string{}, "Define environment variables which should be added when running the command.")
}

func StartStdioTool(nc *nats.Conn, serverName string, cmd string, envs []string, args ...string) error {
	_, err := tool.StartStdioTools(nc, serverName, cmd, envs, args...)
	return err
}
