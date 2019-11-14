package config

import (
	"github.com/spf13/viper"
)

type Client struct {
	Protocol string
	Endpoint string
}

type Debug struct {
	Addr  string
	Token string
	Pprof bool
}

type HTTP struct {
	Addr string
	Root string
}

type GRPC struct {
	Addr string
	Root string
}

type Tracing struct {
	Enabled   bool
	Type      string
	Endpoint  string
	Collector string
	Service   string
}

type Asset struct {
	Path string
}

type Config struct {
	Viper   *viper.Viper
	Client  Client
	Debug   Debug
	HTTP    HTTP
	GRPC    GRPC
	Tracing Tracing
	Asset   Asset
}

func New() *Config {
	return &Config{
		Viper: viper.New(),
	}
}
