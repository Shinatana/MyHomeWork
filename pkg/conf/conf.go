package conf

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

const prefix = "MYAPP"

func NewCfg(configFile string) (*Conf, error) {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFound) {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		return nil, errors.New("config file not found")
	}

	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("dsn", "")
	viper.SetDefault("log_format", "json")
	viper.SetDefault("log_level", "debug")
	viper.SetDefault("http_port", 8080)

	var cfg Conf
	fmt.Println("dsn =", viper.Get("log_level"))
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// validate env
	if cfg.DSN == "" {
		return nil, errors.New("empty DSN")
	}

	if cfg.LogFormat != "json" && cfg.LogFormat != "text" {
		return nil, errors.New("wrong log format, must be json or text")
	}

	if cfg.LogLevel != "debug" && cfg.LogLevel != "info" && cfg.LogLevel != "warn" && cfg.LogLevel != "error" {
		return nil, errors.New("wrong log level, must be debug, info, warn or error")
	}

	if cfg.HttpPort < 1024 || cfg.HttpPort > 65535 {
		return nil, errors.New("port is out of valid range 1024-65535")
	}
	return &cfg, nil
}
