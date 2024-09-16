package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"urlShortener/internal/cfgmodels"
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

	router.Post("/url", save.NewSave(log, storage))
	router.Get("/{alias}", redirect.NewRedirect(log, storage))

	log.Info("start me", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	fmt.Println("I work")
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start srv")
	}
	log.Error("server stopped")

}
