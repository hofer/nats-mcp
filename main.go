package main

import "github.com/hofer/nats-mcp/cmd"

var version string = "0.0.0"

func main() {
	cmd.Execute(version)
}
