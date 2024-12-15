package handler

import (
	"net/http"

	"github.com/rs/cors"
	"github.com/sebest/xff"
	"github.com/tranminhquanq/gomess/internal/app/repository"
	"github.com/tranminhquanq/gomess/internal/app/usecase"
	"github.com/tranminhquanq/gomess/internal/config"
	"github.com/tranminhquanq/gomess/internal/storage"
)

const (
	audHeaderName  = "X-JWT-AUD"
	defaultVersion = "unknown version"
)

type Option interface {
	apply(*APIHandler)
}

type APIHandler struct {
	handler      http.Handler
	db           *storage.Connection
	globalConfig *config.GlobalConfiguration
	version      string
}

func NewHandler(globalConfig *config.GlobalConfiguration, db *storage.Connection, opt ...Option) *APIHandler {
	return NewHandlerWithVersion(globalConfig, db, defaultVersion, opt...)
}

func NewHandlerWithVersion(
	globalConfig *config.GlobalConfiguration,
	db *storage.Connection,
	version string,
	opt ...Option,
) *APIHandler {
	api := &APIHandler{
		globalConfig: globalConfig,
		db:           db,
		version:      version,
	}

	xffmw, _ := xff.Default()

	r := newRouter()

	r.UseBypass(xffmw.Handler)
	r.UseBypass(recoverer)

	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userHandler := NewUserHandler(userUsecase)

	wsHandler := NewWSHandler(userUsecase)

	r.Get("/health", api.HealthCheck)
	r.Get("/ws", wsHandler.HandleWS)

	r.Route("/api", func(r *router) {
		r.Get("/users", userHandler.GetUsers)
		r.Get("/user", userHandler.GetUser)
	})

	corsHandler := cors.New(cors.Options{
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   api.globalConfig.CORS.AllAllowedHeaders([]string{"Accept", "Authorization", "Content-Type", "X-Client-IP", "X-Client-Info", audHeaderName}),
		ExposedHeaders:   []string{"X-Total-Count"},
		AllowCredentials: true,
	})

	api.handler = corsHandler.Handler(r)

	return api
}

// ServeHTTP implements the http.Handler interface by passing the request along
// to its underlying Handler.
func (hdl *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hdl.handler.ServeHTTP(w, r)
}

type HealthCheckResponse struct {
	Version     string `json:"version"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *APIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) error {
	return sendJSON(w, http.StatusOK, HealthCheckResponse{
		Version:     h.version,
		Name:        "GoMess",
		Description: "GoMess is a simple messaging service",
	})
}
