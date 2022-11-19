package v1

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"refactoring/internal/storage"
)

//go:generate mockery --all

var (
	ErrEmptyConfig  = errors.New("config is empty")
	ErrEmptyStorage = errors.New("storage is empty")
)

type API struct {
	server  *http.Server
	storage storage.Storage
}

// New returns new API.
func New(config Config, storage storage.Storage) (newAPI *API, err error) {
	log.Debug().Str("config", config.String()).Msg("api.New START")
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("api.New END")
		} else {
			log.Debug().Msg("api.New END")
		}
	}()

	if config == nil {
		return nil, ErrEmptyConfig
	}
	if storage == nil {
		return nil, ErrEmptyStorage
	}

	newAPI = &API{}

	newAPI.storage = storage

	newServer := &http.Server{}
	newServer.Addr = config.RunAddress()

	router := newAPI.newRouter()
	newServer.Handler = router

	newAPI.server = newServer

	return newAPI, nil
}

func (a *API) newRouter() chi.Router {
	log.Debug().Msg("api.newRouter START")
	defer log.Debug().Msg("api.newRouter END")

	newRouter := chi.NewRouter()

	// TODO
	newRouter.Use(middleware.RequestID) // разобраться, зачем нужен
	newRouter.Use(middleware.RealIP)    // разобраться, зачем нужен
	newRouter.Use(middleware.Logger)
	newRouter.Use(middleware.Recoverer)
	newRouter.Use(middleware.Timeout(60 * time.Second)) // протестировать его работу

	newRouter.Get("/", a.startHandler)

	newRouter.Route("/api", func(r chi.Router) {

		r.Route("/v1", func(r chi.Router) {

			r.Route("/users", func(r chi.Router) {
				r.Get("/", a.searchUsers)
				r.Post("/", a.createUser)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", a.getUser)
					r.Patch("/", a.updateUser)
					r.Delete("/", a.deleteUser)
				})
			})
		})
	})

	return newRouter
}

// Run API starts the API.
func (a *API) Run(ctx context.Context) (err error) {
	log.Debug().Msg("api.New Run")
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("api.Run END")
		} else {
			log.Debug().Msg("api.Run END")
		}
	}()

	defer a.server.Close()
	return a.server.ListenAndServe()
}
