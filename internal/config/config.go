// internal/config/config.go
package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	EthNodeURL      string `yaml:"eth_node_url"`
	EtherscanAPIKey string `yaml:"etherscan_api_key"`
}

const filePath = "C:\\Users\\zmcmanus\\go\\src\\github.com\\zachmdsi\\go-token-cli\\config.yaml"

func LoadConfig() (Config, error) {
	var config Config

	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func SaveConfig(config Config, filePath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, os.FileMode(0644))
}
