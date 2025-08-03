package config

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/spf13/viper"

	"transaction_demo/app/constant"
)

var (
	configFileMap = map[string]string{
		constant.EnvLocal:   "app/config/env/local.yaml",
		constant.EnvDevelop: "app/config/env/develop.yaml",
		constant.EnvProd:    "app/config/env/prod.yaml",
	}
)

// Loader represents the configuration loader for the application
type Loader struct {
	core       *viper.Viper
	env        string
	configFile string
}

// NewLoader creates and initializes a new configuration loader.
// It sets up the environment, determines the config file,
// and initializes the viper core for configuration management.
// Returns a configured Loader instance or an error if initialization fails.
func NewLoader() (*Loader, error) {
	loader := &Loader{}

	env, err := loader.setEnv()
	if err != nil {
		return nil, err
	}
	configFile, err := loader.setConfigFile(env)
	if err != nil {
		return nil, err
	}

	loader.initCore(configFile)

	return loader, nil
}

// initCore initializes the viper core with the specified config file path.
func (l *Loader) initCore(configFilePath string) {
	l.core = viper.New()
	l.core.SetConfigFile(configFilePath)
}

// Load reads the configuration file and executes the provided callback function.
// Returns an error if configuration reading or callback execution fails.
func (l *Loader) Load(do func(core *viper.Viper, env string) error) error {
	l.core.AutomaticEnv()
	l.core.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := l.core.ReadInConfig()
	if err != nil {
		return err
	}

	err = do(l.core, l.env)

	if err != nil {
		return err
	}

	return nil
}

// setEnv determines and sets the application environment.
// Returns the environment string or an error if the environment is invalid.
func (l *Loader) setEnv() (string, error) {
	envVal := os.Getenv("APP_ENV")
	if envVal == "" {
		envVal = constant.EnvLocal
	}
	if !validEnv(envVal) {
		return "", fmt.Errorf("invalid env: %v", envVal)
	}
	l.env = envVal
	return l.env, nil
}

// validEnv checks if the provided environment string is valid.
// Returns true if the environment is valid, false otherwise.
func validEnv(env string) bool {
	envList := []string{constant.EnvLocal, constant.EnvDevelop, constant.EnvProd}
	return slices.Contains(envList, env)
}

// setConfigFile determines the configuration file path based on the environment.
// Returns the config file path or an error if no config file is found for the environment.
func (l *Loader) setConfigFile(env string) (string, error) {
	configFile, ok := configFileMap[env]
	if !ok {
		return "", fmt.Errorf("config file not found for env: %s", env)
	}

	l.configFile = configFile
	return configFile, nil
}
