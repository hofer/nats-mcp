package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A collection of commands which show details of exposed tools",
	Long: `The client sub commands can be used to list details of all MCP tools which are exposed via NATS.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("admin called")
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
