package cmd

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tranminhquanq/gomess/internal/app/handler"
	"github.com/tranminhquanq/gomess/internal/config"
)

const (
	SlowlorisTimeout = 2 * time.Second // Time to mitigate Slowloris attack
)

var serveCmd = cobra.Command{
	Use:  "serve",
	Long: "Start API server",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd.Context())
	},
}

func serve(ctx context.Context) {
	if err := config.LoadFile((configFile)); err != nil {
		logrus.WithError(err).Fatal("unable to load config")
	}

	config, err := config.LoadGlobalFromEnv()
	if err != nil {
		logrus.WithError(err).Fatal("unable to load config")
	}

	addr := net.JoinHostPort(config.API.Host, config.API.Port)

	baseCtx, baseCancel := context.WithCancel(context.Background())
	defer baseCancel()

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           handler.NewHandler(config, "1.0.0"), // Version is hardcoded for now
		ReadHeaderTimeout: SlowlorisTimeout,
		BaseContext: func(l net.Listener) context.Context {
			return baseCtx
		},
	}

	log := logrus.WithField("component", "api")

	var wg sync.WaitGroup
	defer wg.Wait() // Do not return to caller until this goroutine is done.

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()

		defer baseCancel() // close baseContext

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Minute)
		defer shutdownCancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
			log.WithError(err).Error("shutdown failed")
		}
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.WithError(err).Fatal("http server listen failed")
	}
}
