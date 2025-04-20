package cmd

import (
	"github.com/hofer/nats-mcp/internal/client"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all MCP tools accessible via NATS",
	Long:  `Searches for all tools which are exposed via NATS and prints a list.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Listing all tools...")
		nc, err := nats.Connect(natsUrl)
		if err != nil {
			log.Fatal(err)
		}
		err = client.ListTools(nc)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	clientCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&natsUrl, "url", "u", os.Getenv("NATS_URL"), "URL to the Nats.io server")
	if os.Getenv("NATS_URL") == "" {
		listCmd.MarkFlagRequired("url")
	}
}
