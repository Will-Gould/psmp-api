package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
	"github.com/Will-Gould/psmp-api/internal/grieflogger"
	"github.com/Will-Gould/psmp-api/internal/players"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type application struct {
	config config
	// logger
	db *sql.DB
}

// mount
func (app application) mount() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set timeout
	r.Use(middleware.Timeout(60 * time.Second))

	// New repo
	repo := repo.New(app.db)

	// start grief logger service & handler
	griefloggerService := grieflogger.NewService(repo)
	griefLoggerHandler := grieflogger.NewHandler(griefloggerService)

	// start player service & handler
	playerService := players.NewService(repo)
	playerHandler := players.NewHandler(playerService, griefLoggerHandler)
	r.Get("/players", playerHandler.ListPlayersHandler)
	r.Get("/players/{uuid}", playerHandler.ListPlayer)
	r.Get("/players/{uuid}/overview", playerHandler.GetPlayerOverview)
	r.Get("/players/{uuid}/block-data", playerHandler.ListPlayerBlockData)

	return r
}

// run
func (app application) run(h http.Handler) error {
	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	message := "Server has started at address:" + app.config.addr
	slog.Info(message)

	return srv.ListenAndServe()
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
