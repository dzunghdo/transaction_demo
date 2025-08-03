package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var configInstance *Config

// InitConfig initializes the configuration for the application
// by loading it from the environment or configuration files using viper.
func InitConfig() (*Config, error) {
	if configInstance != nil {
		return configInstance, nil
	}

	loader, err := NewLoader()
	if err != nil {
		return nil, err
	}

	c := &Config{}

	err = loader.Load(func(core *viper.Viper, env string) error {
		err := core.Unmarshal(c)
		if err != nil {
			return err
		}
		c.Env = env
		return nil
	})
	if err != nil {
		return nil, err
	}

	configInstance = c

	if c.Server.Port == 0 {
		return nil, fmt.Errorf("invalid server port: %v", c.Server.Port)
	}

	return configInstance, nil
}
