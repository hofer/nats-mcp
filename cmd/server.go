package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/hofer/nats-mcp/internal/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start an MCP Server exposing all NATS MCP tools",
	Long: `This command will start a MCP Server and make available all tools which are defined
as NATS microservices as tools.

The server is accessible via stdio only. 
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")

		natsUrl, _ := cmd.Flags().GetString("url")
		//natsUrl := os.Getenv("NATS_SERVER_URL")
		server.StartServer(natsUrl)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().String("url", "", "URL to the Nats.io server")
	serverCmd.MarkFlagRequired("url")
}
