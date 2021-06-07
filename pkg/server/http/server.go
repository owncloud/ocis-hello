package http

import (
	"encoding/json"
	"net/http"

	"github.com/asim/go-micro/v3"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/owncloud/ocis-hello/pkg/assets"
	"github.com/owncloud/ocis-hello/pkg/proto/v0"
	"github.com/owncloud/ocis-hello/pkg/version"
	"github.com/owncloud/ocis/ocis-pkg/account"
	"github.com/owncloud/ocis/ocis-pkg/middleware"
	ohttp "github.com/owncloud/ocis/ocis-pkg/service/http"
)

type greetRequest struct {
	Name string `json:"name"`
}

// Server initializes the http service and server.
func Server(opts ...Option) ohttp.Service {
	options := newOptions(opts...)
	handler := options.Handler

	svc := ohttp.NewService(
		ohttp.Logger(options.Logger),
		ohttp.Name(options.Name),
		ohttp.Version(options.Config.Server.Version),
		ohttp.Address(options.Config.HTTP.Addr),
		ohttp.Namespace(options.Config.HTTP.Namespace),
		ohttp.Context(options.Context),
		ohttp.Flags(options.Flags...),
	)

	mux := chi.NewMux()

	mux.Use(middleware.RealIP)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.NoCache)
	mux.Use(middleware.Cors)
	mux.Use(middleware.Secure)
	mux.Use(middleware.ExtractAccountUUID(
		account.Logger(options.Logger),
		account.JWTSecret(options.Config.TokenManager.JWTSecret)),
	)

	mux.Use(middleware.Version(
		options.Name,
		version.String,
	))

	mux.Use(middleware.Logger(
		options.Logger,
	))

	mux.Use(middleware.Static(
		options.Config.HTTP.Root,
		assets.New(
			assets.Logger(options.Logger),
			assets.Config(options.Config),
		),
		options.Config.HTTP.CacheTTL,
	))

	mux.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Post("/api/v0/greet", func(w http.ResponseWriter, r *http.Request) {
			var req greetRequest

			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if req.Name == "" {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "missing a name")
				return
			}

			var accountID string
			val := r.Context().Value(middleware.UUIDKey)
			if val != nil {
				accountID = val.(string)
			}
			greeting := handler.Greet(accountID, req.Name)

			rsp := &proto.GreetResponse{
				Message: greeting,
			}

			render.Status(r, http.StatusCreated)
			render.JSON(w, r, rsp)
		})
	})

	err := micro.RegisterHandler(svc.Server(), mux)
	if err != nil {
		options.Logger.Fatal().Err(err).Msg("failed to register the handler")
	}

	svc.Init()
	return svc
}
