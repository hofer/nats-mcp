package cmd

import (
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

var configFile string

var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Expose tools from a local MCP Server via NATS",
	Long: `This command can be used to expose local MCP tools (a MCP server started locally) via NATS. With
just a few simple commands many different MCP servers can be made accessible via NATS.
`,
	Run: func(cmd *cobra.Command, args []string) {
		nc, err := nats.Connect(natsUrl)
		if err != nil {
			log.Fatal(err)
		}

		config, err := LoadConfig(configFile)
		if err != nil {
			log.Fatal(err)
		}

		// Start all Stdio tools defined in the config file:
		for sName, c := range config.GetStdioServers() {
			// TODO: Fix passing envs...
			log.Infof("Starting Stdio tool '%s'", sName)
			err = StartStdioTool(nc, sName, c.Command, []string{}, c.Args...)
			if err != nil {
				log.Error(err)
			}
		}

		// Start all Stdio tools defined in the config file:
		for sName, c := range config.GetSseServers() {
			// TODO: Fix passing envs...
			log.Infof("Starting SSE tool '%s'", sName)
			err = StartSseTool(nc, sName, c.Url)
			if err != nil {
				log.Error(err)
			}
		}

		log.Info("Waiting for incoming tool calls...")
		runtime.Goexit()
	},
}

func init() {
	rootCmd.AddCommand(toolCmd)
	toolCmd.Flags().StringVarP(&natsUrl, "url", "u", os.Getenv("NATS_URL"), "URL to the Nats.io server")
	if os.Getenv("NATS_URL") == "" {
		toolCmd.MarkFlagRequired("url")
	}

	toolCmd.Flags().StringVarP(&configFile, "file", "f", "", "JSON config file containing MCP server configurations.")
	toolCmd.MarkFlagRequired("file")
}
