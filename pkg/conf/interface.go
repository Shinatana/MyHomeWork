package conf

type Conf struct {
	DSN       string `mapstructure:"dsn"`
	LogFormat string `mapstructure:"log_format"`
	LogLevel  string `mapstructure:"log_level"`
	HttpPort  int    `mapstructure:"http_port"`
}
