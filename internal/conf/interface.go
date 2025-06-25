package conf

type Conf struct {
	DSN       string `mapstructure:"dsn" validate:"required, uri"`
	LogFormat string `mapstructure:"log_format" validate:"oneof=json text"`
	LogLevel  string `mapstructure:"log_level" validate:"oneof=debug info warn error"`
	HttpPort  int    `mapstructure:"http_port" validate:"gte=1024,lte=65535"`
}
