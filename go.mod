module github.com/owncloud/ocis-hello

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/cespare/reflex v0.2.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/grpc-gateway v1.14.6
	github.com/haya14busa/goverage v0.0.0-20180129164344-eec3514a20b5
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/kr/pty v1.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/mitchellh/gox v1.0.1
	github.com/ogier/pflag v0.0.1 // indirect
	github.com/oklog/run v1.1.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis-pkg/v2 v2.4.0
	github.com/owncloud/ocis-settings v0.3.2-0.20200827193534-8caf098e6537
	github.com/prometheus/client_golang v1.7.1
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.6.1
	go.opencensus.io v0.22.4
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b
	golang.org/x/net v0.0.0-20200625001655-4c5254603344
	google.golang.org/genproto v0.0.0-20200513103714-09dca8ec2884
	google.golang.org/protobuf v1.23.0
	honnef.co/go/tools v0.0.1-2020.1.5
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
