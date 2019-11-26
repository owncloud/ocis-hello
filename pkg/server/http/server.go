package http

import (
	"fmt"
	"time"

	"github.com/micro/go-micro/web"
	v0proto "github.com/owncloud/ocis-hello/pkg/proto/v0"
	"github.com/owncloud/ocis-hello/pkg/version"
)

func Server(opts ...Option) (web.Service, error) {
	options := newOptions(opts...)

	options.Logger.Debug().Msg("Server!")

	service := web.NewService(
		web.Name("go.micro.web.hello"),
		web.Version(version.String),
		web.RegisterTTL(time.Second*30),
		web.RegisterInterval(time.Second*10),
	)

	options.Logger.Debug().Msg("After service...")

	v0proto.RegisterHelloWeb(
		service,
		"go.micro.api.hello",
		nil,
	)

	options.Logger.Debug().Msg("After register...")

	fmt.Printf("%+v\n", service.Options().Service.Options().Cmd)

	if err := service.Init(); err != nil {

		options.Logger.Debug().Msg("Init failed...")

		return nil, err
	}

	options.Logger.Debug().Msg("After init...")

	return service, nil
}
