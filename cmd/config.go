package cmd

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
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
	Servers map[string]interface{} `json:"servers"`
}

func (c McpConfig) GetStdioServers() map[string]McpServerStdioConfig {
	stdioServers := map[string]McpServerStdioConfig{}
	for sName, server := range c.Servers {
		if serverConfig, ok := server.(map[string]interface{}); ok {
			if serverConfig["type"] == "stdio" {
				var u McpServerStdioConfig
				err := decodeConfig(&u, serverConfig)
				if err != nil {
					log.Warnf("Error decoding serverConfig: %v", err)
					continue
				}
				stdioServers[sName] = u
			}
		}
	}
	return stdioServers
}

func (c McpConfig) GetSseServers() map[string]McpServerSseConfig {
	sseServers := map[string]McpServerSseConfig{}
	for sName, server := range c.Servers {
		if serverConfig, ok := server.(map[string]interface{}); ok {
			if serverConfig["type"] == "sse" {
				var u McpServerSseConfig
				err := decodeConfig(&u, serverConfig)
				if err != nil {
					log.Warnf("Error decoding serverConfig: %v", err)
					continue
				}
				sseServers[sName] = u
			}
		}
	}
	return sseServers
}

func decodeConfig(u interface{}, serverConfig map[string]interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   u,
		TagName:  "json", // Specify the tag name to use
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	err = decoder.Decode(serverConfig)
	return err
}

type McpServerConfig struct {
	Type string `json:"type"`
}

type McpServerStdioConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

type McpServerSseConfig struct {
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}
