package cmd

import (
	"encoding/json"
	"os"
)

func LoadConfig(configFile string) (McpConfig, error) {
	b, err := os.ReadFile(configFile)
	if err != nil {
		return McpConfig{}, err
	}

	var config McpConfig
	err = json.Unmarshal(b, &config)
	return config, err
}

type McpConfig struct {
	Servers map[string]McpServerConfig `json:"servers"`
}

type McpServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}
