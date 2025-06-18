package conf

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Conf struct {
	DbDSN               string
	LogFormat, LogLevel string
	HttpPort            int
}

const prefix = "MYAPP_"

func NewCfg() (*Conf, error) {
	dsn, ok := os.LookupEnv(prefix + "DB_DSN")
	if !ok {
		return nil, fmt.Errorf("env %s not found", prefix+"DB_DSN")
	}

	logFormat, ok := os.LookupEnv(prefix + "LOG_FORMAT")
	if !ok {
		return nil, fmt.Errorf("env %s not found", prefix+"LOG_FORMAT")
	}

	logLevel, ok := os.LookupEnv(prefix + "LOG_LEVEL")
	if !ok {
		return nil, fmt.Errorf("env %s not found", prefix+"LOG_LEVEL")
	}

	port, ok := os.LookupEnv(prefix + "HTTP_PORT")
	if !ok {
		return nil, fmt.Errorf("env %s not found", prefix+"HTTP_PORT")
	}
	httpPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("env %s is not a valid port", prefix+"HTTP_PORT")
	}

	// validate env
	if dsn == "" {
		return nil, errors.New("empty dsn")
	}

	if logFormat != "json" && logFormat != "text" {
		return nil, errors.New("wrong log format, must be json or text")
	}

	if logLevel != "debug" && logLevel != "info" && logLevel != "warn" && logLevel != "error" {
		return nil, errors.New("wrong log level, must be debug, info, warn or error")
	}

	if httpPort < 1024 || httpPort > 65535 {
		return nil, errors.New("port is out of valid range 1024-65535")
	}

	return &Conf{DbDSN: dsn,
			LogFormat: logFormat,
			LogLevel:  logLevel,
			HttpPort:  httpPort,
		},
		nil
}
