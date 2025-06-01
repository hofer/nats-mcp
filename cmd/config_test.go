package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	c, err := LoadConfig("../example-mcp.json")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Len(t, c.Servers, 1)
}
