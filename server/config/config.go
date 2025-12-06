package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port 			int `yaml:"port"`
		ReadTimeoutMS 	int `yaml:"read_timeout_ms"`
		WriteTimeoutMS 	int `yaml:"write_timeout_ms"`
		MaxConnections 	int `yaml:"max_connections"`
		Workers  		int `yaml:"workers"`
	} `yaml:"server"`

	Logging struct {
		Level 			string `yaml:"level"`
		Format 			string `yaml:"format"`
		DisableReqLogs 	bool `yaml:"disable_request_logs"`
	} `yaml:"logging"`

	Metrics struct {
		Enabled 		bool `yaml:"enabled"`
		Port    		int  `yaml:"port"`
		Path   			string `yaml:"path"`
	} `yaml:"metrics"`

	Performance struct {
		ParseJSON 			bool `yaml:"parse_json"`
		IncludeTimestamp 	bool `yaml:"include_timestamp"`
		MinimizeResponse 	bool `yaml:"minimize_response"`
	} `yaml:"performance"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}