package cmd

import (
	"github.com/spf13/cobra"
	//"github.com/hofer/nats-mcp/internal/client"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all MCP tools accessible via NATS",
	Long: `Searches for all tools which are exposed via NATS and prints a list.`,
	Run: func(cmd *cobra.Command, args []string) {
		//client.StartClient()
	},
}

func init() {
	clientCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
