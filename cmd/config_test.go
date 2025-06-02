package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	c, err := LoadConfig("../example-mcp.json")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Len(t, c.Servers, 2)

	// Check if the stdio server is correctly parsed
	stdioServers := c.GetStdioServers()
	assert.Len(t, stdioServers, 1)

	// Check if the sse server is correctly parsed
	sseServers := c.GetSseServers()
	assert.Len(t, sseServers, 1)
}
