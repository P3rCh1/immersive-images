package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/P3rCh1/immersive-images/backend/internal/api/handlers"
	"github.com/P3rCh1/immersive-images/backend/internal/api/server"
	"github.com/P3rCh1/immersive-images/backend/internal/config"
	"github.com/P3rCh1/immersive-images/backend/internal/dal/object"
	"github.com/P3rCh1/immersive-images/backend/internal/dal/postgres"
	"github.com/P3rCh1/immersive-images/backend/internal/logger"
	"github.com/go-chi/chi/v5"
)

func main() {
	configPath := flag.String("c", "/etc/config.yaml", "specify configuration file")
	flag.Parse()

	if err := config.InitConfig(*configPath); err != nil {
		fmt.Fprintf(os.Stderr, "failed to init config: %v", err)
		os.Exit(1)
	}

	log, err := logger.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logger: %v", err)
		os.Exit(1)
	}

	if err := run(log); err != nil {
		log.Error(
			"execute service error",
			"error", err,
		)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	db, err := postgres.New(log)
	if err != nil {
		return err
	}
	defer db.Close()

	s3 := object.New(log)

	apiHandlers := handlers.New(log, db, s3)

	router := chi.NewRouter()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	httpServer := http.Server{
		Addr: config.Config.Server.Host + ":" + config.Config.Server.Port,
		Handler: server.HandlerWithOptions(apiHandlers, server.ChiServerOptions{
			BaseURL:          "/v1",
			BaseRouter:       router,
			Middlewares:      []server.MiddlewareFunc{apiHandlers.Recover()},
			ErrorHandlerFunc: apiHandlers.ValidationErrorHandler,
		}),
		ReadTimeout:  config.Config.Server.ReadTimeout,
		WriteTimeout: config.Config.Server.WriteTimeout,
		IdleTimeout:  config.Config.Server.IdleTimeout,
	}

	go func() {
		<-ctx.Done()

		sdCtx, sdCancel := context.WithTimeout(context.Background(), config.Config.Server.ShutdownTimeout)
		defer sdCancel()

		if err := httpServer.Shutdown(sdCtx); err != nil {
			log.Error(
				"failed to shutdown server",
				"error", err,
			)
		}
	}()

	log.Info("start server")

	if err := httpServer.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server crashed: %w", err)
		}

		log.Info("server stopped")
	}

	return nil
}
