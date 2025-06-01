package natsmcp

//func StartTool(nc *nats.Conn, transport *Nats) error {
//	log.Infof("Starting McpEchoService...")
//
//	// Create tools and corresponding Handlers:
//	tools := NatsMcpTool{
//		Tool: mcp.NewTool("hello_echo",
//			mcp.WithDescription("Say hello to someone"),
//			mcp.WithString("name",
//				mcp.Required(),
//				mcp.Description("Name of the person to greet"),
//			)),
//		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
//			log.Infof("Tool call for hello_echo.")
//			name, ok := request.GetArguments()["name"].(string)
//			if !ok {
//				return mcp.NewToolResultError("name must be a string"), nil
//			}
//			return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
//		},
//	}
//
//	serviceName := "EchoService"
//
//	// Add all tools to a toolbox
//	toolBox := NewNatsMcpToolBox(tools)
//
//	// Expose all tools in the toolbox as a Nats microservice
//	srv, err := toolBox.AddToolsAsNatsService(nc, serviceName)
//	if err != nil {
//		return err
//	}
//
//	toolBox.AddTransportNatsService(srv, transport, serviceName)
//
//	log.Infof("Service connected at %s", nc.ConnectedAddr())
//	return nil
//}
