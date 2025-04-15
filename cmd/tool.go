package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/hofer/nats-mcp/internal/tool"
)

var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Expose tools from a MCP Server via NATS",
	Long: `This command can be used to expose local MCP tools (a MCP server started locally) via NATS. With
just a few simple commands many other MCP servers can be made accessible via NATS.

`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tool called")
		fmt.Println(args)

		natsUrl, _ := cmd.Flags().GetString("url")
		command, _ := cmd.Flags().GetString("command")
		tool.StartTool(natsUrl, command, args)
	},
}

func init() {
	rootCmd.AddCommand(toolCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// toolCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// toolCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	toolCmd.Flags().String("url", "", "URL to the Nats.io server")
	toolCmd.MarkFlagRequired("url")
	toolCmd.Flags().String("command", "", "Help message for toggle")
	toolCmd.MarkFlagRequired("command")
}
