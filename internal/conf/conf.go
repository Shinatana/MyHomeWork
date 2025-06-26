package conf

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const (
	prefix = "MYAPP"
)

func NewCfg(configFile string) (*Conf, error) {

	viper.SetConfigFile(configFile)

	_ = viper.ReadInConfig()

	viper.SetEnvPrefix(prefix)
	viper.AutomaticEnv()

	var cfg Conf

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("failed validation: %w", err)
	}

	return &cfg, nil
}
