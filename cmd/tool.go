package cmd

import (
	"github.com/hofer/nats-mcp/internal/tool"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

var configFile string
var commandStr string
var toolServerName string
var environment []string

var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Expose tools from a local MCP Server (Stdio) via NATS",
	Long: `This command can be used to expose local MCP tools (a MCP server started locally) via NATS. With
just a few simple commands many different MCP servers can be made accessible via NATS.
`,
	Run: func(cmd *cobra.Command, args []string) {
		nc, err := nats.Connect(natsUrl)
		if err != nil {
			log.Fatal(err)
		}

		if len(configFile) != 0 {
			config, err := LoadConfig(configFile)
			if err != nil {
				log.Fatal(err)
			}

			for sName, c := range config.Servers {
				// TODO: Fix passing envs...
				log.Infof("Starting tool '%s'", sName)
				err = StartTool(nc, sName, c.Command, []string{}, c.Args...)
				if err != nil {
					log.Fatal(err)
				}
			}

		} else {
			err = StartTool(nc, toolServerName, commandStr, environment, args...)
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Info("Waiting for incoming tool calls...")
		runtime.Goexit()
	},
}

func StartTool(nc *nats.Conn, serverName string, cmd string, envs []string, args ...string) error {
	_, err := tool.StartTool(nc, serverName, cmd, envs, args...)
	return err
}

func init() {
	rootCmd.AddCommand(toolCmd)
	toolCmd.Flags().StringVarP(&natsUrl, "url", "u", os.Getenv("NATS_URL"), "URL to the Nats.io server")
	if os.Getenv("NATS_URL") == "" {
		toolCmd.MarkFlagRequired("url")
	}

	toolCmd.Flags().StringVarP(&toolServerName, "name", "n", "", "Server name (if used with commandline args)")
	//	toolCmd.MarkFlagRequired("name")

	toolCmd.Flags().StringVarP(&configFile, "file", "f", "", "JSON config file containing MCP server configurations.")

	toolCmd.Flags().StringVarP(&commandStr, "command", "c", "", "Command to start the local MCP Server")
	//toolCmd.MarkFlagRequired("command")

	toolCmd.Flags().StringArrayVarP(&environment, "env", "e", []string{}, "Define environment variables which should be added when running the command.")
}
