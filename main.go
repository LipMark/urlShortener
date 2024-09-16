package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"urlShortener/internal/cfgmodels"
	"urlShortener/internal/http-server/handlers/deletehandler"
	"urlShortener/internal/http-server/handlers/redirect"
	"urlShortener/internal/http-server/handlers/save"
	"urlShortener/internal/storage/sqlite"
	"urlShortener/util"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

// loads values from .env into the system
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	cfg := cfgmodels.ConfigSetUp()

	// SetUp logger
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	storage, err := sqlite.CreateStorage(cfg.StoragePath)

	if err != nil {
		log.Error("failed to init storage", util.SlogErr(err))
		os.Exit(1)
	}
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("urlShortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.NewSave(log, storage))
		r.Delete("/{id}", deletehandler.NewDelete(log, storage))
	})

	router.Post("/url", save.NewSave(log, storage))
	router.Get("/{alias}", redirect.NewRedirect(log, storage))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		fmt.Println("Server started on port", 8080)
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("listen:%s\n", err)
		}
	}()
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", util.SlogErr(err))

		return
	}

	fmt.Println("server stopped")
}
