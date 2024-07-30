package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	WorkerInterval string   `mapstructure:"workerInterval"`
	ValidTokens    []string `mapstructure:"validTokens"`
	NumUsers       int      `mapstructure:"numUsers"`
	NumMsg         int      `mapstructure:"numMsg"`
	CountWorkers   int      `mapstructure:"countWorkers"`
}

// GetConfig load project settings from .yaml file
func GetConfig(configFile string) (*Config, error) {
	var config Config
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return &config, nil
}
