module nats-mcp/mcp-example

go 1.24.1

require (
	github.com/hofer/nats-mcp v0.0.0
	github.com/mark3labs/mcp-go v0.20.1
	github.com/nats-io/nats.go v1.41.1
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/nats-io/nkeys v0.4.9 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)

replace github.com/hofer/nats-mcp v0.0.0 => ../../
