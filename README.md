# 🧰 NATS MCP toolbox
![workflow](https://github.com/hofer/nats-mcp/actions/workflows/build.yml/badge.svg)

This is a collection of cli tools to expose MCP tools via NATS microservices. It can be used either standalone (exposing
local/existing MCP Servers via NATS) or as library in Go to expose functionality as tools.

> [!WARNING]
> 🚨 🚧 This tool is under active development 🚧 🚨
>
> This tool is very much work in progress. While the tools should all work, expect almost
> daily breaking changes.

## Usage

Exposing an existing MCP Server via NATS:
```
./nats-mcp tool --url "nats://localhost:4222" --command="./whatever-mcp-server" arg1 arg2
```

To check what MCP tools are exposed via NATS use the following command:
```
./nats-mcp client list --url "nats://localhost:4222"
```

To then actually use the tool with your local agent, add the following command to your MCP config:
```
./nats-mcp server --url "nats://localhost:4222"
```

Given the `--url` parameter is used for allmost all arguments, the `NATS_URL` environment variable is used as a default.



