package conf

type Conf struct {
	DSN       string `mapstructure:"dsn" validate:"required"`
	LogFormat string `mapstructure:"log_format" validate:"oneof=json text"`
	LogLevel  string `mapstructure:"log_level" validate:"oneof=debug info warn error"`
	HttpPort  int    `mapstructure:"http_port" validate:"gt=1024,lt=65535"`
}
