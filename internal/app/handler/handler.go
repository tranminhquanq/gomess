package handler

import (
	"net/http"

	"github.com/tranminhquanq/gomess/internal/config"
)

func NewHandler(config *config.GlobalConfiguration, version string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello, world!"))
	})
}
