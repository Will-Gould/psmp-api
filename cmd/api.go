package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
	"github.com/Will-Gould/psmp-api/internal/mapping"
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

	// Start mapping service & handler
	mappingService := mapping.NewService(repo)
	mappingHandler := mapping.NewHandler(mappingService)
	// load mapping data
	materials := mappingHandler.GetMaterials(context.Background())
	bannedPlacedMaterials, bannedBrokenMaterials := mappingHandler.GetBannedMaterials(context.Background(), materials)
	causeMapping := mappingHandler.GetDeathCauseMapping(context.Background())
	mobMapping := mappingHandler.GetMobMapping(context.Background())
	advancementMapping := mappingHandler.GetAdvancementMapping(context.Background())

	mappingData := &mapping.MappingData{
		Materials:             materials,
		BannedPlacedMaterials: bannedPlacedMaterials,
		BannedBrokenMaterials: bannedBrokenMaterials,
		CauseMapping:          causeMapping,
		MobMapping:            mobMapping,
		AdvancementMapping:    advancementMapping,
	}

	// schedule mapping data updates

	// Start player service & initialise in memory leaderboards
	playerService := players.NewService(repo)
	playerHandler := players.NewHandler(playerService)
	playerHandler.Initialise(context.Background(), mappingData)

	// schedule player list & leaderboard updates

	r.Get("/api/players/{uuid}", playerHandler.GetServerPlayer)
	r.Get("/api/players/leaderboards/server", playerHandler.ListServerLeaderboard)
	r.Get("/api/players/leaderboards/combat", playerHandler.ListCombatLeaderboard)

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

	message := "Server has started on port" + app.config.addr
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
