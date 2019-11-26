package config

type Log struct {
	Level  string
	Pretty bool
	Color  bool
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
	File    string
	Log     Log
	Debug   Debug
	HTTP    HTTP
	GRPC    GRPC
	Tracing Tracing
	Asset   Asset
}

func New() *Config {
	return &Config{}
}
