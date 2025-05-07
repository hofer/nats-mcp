# 🧰 NATS MCP toolbox
![workflow](https://github.com/hofer/nats-mcp/actions/workflows/build.yml/badge.svg)

This is a collection of cli tools to expose MCP tools via NATS microservices. It can be used either standalone (exposing
local/existing MCP Servers via NATS) or as library in Go to expose functionality as tools.

> [!WARNING]
> 🚨 🚧 This tool is under active development 🚧 🚨
>
> This tool is very much work in progress. While the tools should all work, expect almost
> daily breaking changes. Please also keep in mind that we currently support stdio servers only.

## Usage

Exposing an existing MCP Server via NATS:
```
./nats-mcp tool --url "nats://localhost:4222" --command="./whatever-mcp-erver" -arg foo
```

To check what MCP tools are exposed via NATS use the following command:
```
./nats-mcp client list --url "nats://localhost:4222"
```

To then actually use the tool with your local agent, add the following command to your MCP config:
```
./nats-mcp server --url "nats://localhost:4222"
```


## Useful links / inspiration:

The following collection of links have been an inspiration for this project.

### MCP Spec:
- https://modelcontextprotocol.io/docs/concepts/tools

### MCP SDKs:
- https://github.com/metoro-io/mcp-golang or https://mcpgolang.com/introduction
- https://github.com/mark3labs/mcp-go
- https://github.com/mark3labs/mcphost

### NATS:
- Nats Service Discovery (code example): https://github.com/nats-io/natscli/blob/main/cli/service_command.go
- Nats ADR: https://github.com/nats-io/nats-architecture-and-design
- Nats Micro: https://github.com/nats-io/nats.go/blob/main/micro/README.md