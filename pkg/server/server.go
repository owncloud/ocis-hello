package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus"
	"github.com/micro/go-plugins/wrapper/trace/opencensus"
	"github.com/oklog/run"
	v0proto "github.com/owncloud/ocis-hello/pkg/proto/v0"
	v0svc "github.com/owncloud/ocis-hello/pkg/service/v0"
	"github.com/owncloud/ocis-hello/pkg/version"
)

func main() {
	var gr run.Group

	grpc := micro.NewService(
		micro.Name("go.micro.api.hello"),
		micro.Version(version.String),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		micro.WrapClient(opencensus.NewClientWrapper()),
		micro.WrapHandler(opencensus.NewHandlerWrapper()),
		micro.WrapSubscriber(opencensus.NewSubscriberWrapper()),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)

	v0proto.RegisterHelloHandler(
		grpc.Server(),
		&v0svc.Hello{},
	)

	gr.Add(func() error {
		grpc.Init()

		if err := grpc.Run(); err != nil {
			return err
		}

		return nil
	}, func(reason error) {
		log.Println(reason)
	})

	http := web.NewService(
		web.Name("go.micro.web.hello"),
		web.Version(version.String),
		web.RegisterTTL(time.Second*30),
		web.RegisterInterval(time.Second*10),
	)

	v0proto.RegisterHelloWeb(
		http,
		grpc.Name(),
		grpc.Client(),
	)

	gr.Add(func() error {
		if err := http.Init(); err != nil {
			return err
		}

		if err := http.Run(); err != nil {
			return err
		}

		return nil
	}, func(reason error) {
		log.Println(reason)
	})

	stop := make(chan os.Signal, 1)

	gr.Add(func() error {
		signal.Notify(stop, os.Interrupt)

		<-stop

		return nil
	}, func(err error) {
		close(stop)
	})

	if err := gr.Run(); err != nil {
		log.Fatal(err)
	}
}
