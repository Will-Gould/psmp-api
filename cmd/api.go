package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
	"github.com/Will-Gould/psmp-api/internal/cache"
	"github.com/Will-Gould/psmp-api/internal/mapping"
	"github.com/Will-Gould/psmp-api/internal/players"
	"github.com/Will-Gould/psmp-api/internal/profiles"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
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

	// initialise cache
	dataStore := cache.DataStore{
		Players:             make(map[string]responsemodels.ServerPlayer),
		MappingData:         cache.MappingData{},
		CombatLeaderboard:   make(map[string]responsemodels.CombatLeaderboardPlayer),
		CraftingLeaderboard: make(map[string]responsemodels.CraftingLeaderboardPlayer),
		StoryLeaderboard:    make(map[string]responsemodels.StoryLeaderboardPlayer),
		StatLeaderboards: cache.StatLeaderboards{
			BlocksPlacedLeaderboard:  map[string]responsemodels.LeaderboardPlayer{},
			BlocksBrokenLeaderboard:  map[string]responsemodels.LeaderboardPlayer{},
			DiamondsMinedLeaderboard: map[string]responsemodels.LeaderboardPlayer{},
			TimePlayedLeaderboard:    map[string]responsemodels.LeaderboardPlayer{},
			PvpKillsLeaderboard:      map[string]responsemodels.LeaderboardPlayer{},
			DeathsLeaderboard:        map[string]responsemodels.LeaderboardPlayer{},
			MobKillsLeaderboard:      map[string]responsemodels.LeaderboardPlayer{},
			PvpKdRatioLeaderboard:    map[string]responsemodels.LeaderboardPlayer{},
		},
	}

	// New repo
	repo := repo.New(app.db)

	// Start mapping service & handler
	mappingService := mapping.NewService(repo)
	mappingHandler := mapping.NewHandler(mappingService, &dataStore)
	// Start player service & handler
	playerService := players.NewService(repo)
	playerHandler := players.NewHandler(playerService, &dataStore)

	// Load data into cache
	mappingHandler.LoadMappingData(context.Background())
	playerHandler.Load(context.Background())

	// create go routine for cache updates
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for range ticker.C {
			slog.Log(context.Background(), slog.LevelInfo, "Updating cache...")
			dataStore.Mu.Lock()
			mappingHandler.LoadMappingData(context.Background())
			playerHandler.Load(context.Background())
			dataStore.Mu.Unlock()
		}
	}()

	// Start profile service & handler
	profileService := profiles.NewService(repo)
	profileHandler := profiles.NewHandler(profileService, &dataStore)

	// Map endpoints
	// players
	r.Get("/api/players/{uuid}", playerHandler.GetServerPlayer)
	r.Get("/api/players/leaderboards/server", playerHandler.ListServerLeaderboard)
	r.Get("/api/players/leaderboards/server/top-ten", playerHandler.ListServerTopTen)
	r.Get("/api/players/leaderboards/combat", playerHandler.ListCombatLeaderboard)
	r.Get("/api/players/leaderboards/combat/top-ten", playerHandler.ListCombatTopTen)
	r.Get("/api/players/leaderboards/crafting", playerHandler.ListCraftingLeaderboard)
	r.Get("/api/players/leaderboards/crafting/top-ten", playerHandler.ListCraftingTopTen)
	r.Get("/api/players/leaderboards/story", playerHandler.ListStoryLeaderboard)
	r.Get("/api/players/leaderboards/combat/{uuid}", playerHandler.GetCombatLeaderboardPlayer)
	r.Get("/api/players/leaderboards/crafting/{uuid}", playerHandler.GetCraftingLeaderboardPlayer)
	r.Get("/api/players/leaderboards/story/{uuid}", playerHandler.GetStoryLeaderboardPlayer)

	// stat leaderboards
	r.Get("/api/players/leaderboards/stats/{uuid}", playerHandler.GetStatRanks)
	r.Get("/api/players/leaderboards/stats/stat/{stat}", playerHandler.ListStatLeaderboard)

	// features
	r.Get("/api/players/leaderboards/champion-player", playerHandler.GetChampionPlayer)
	r.Get("/api/players/leaderboards/most-dangerous-player", playerHandler.GetMostDangerousPlayer)
	r.Get("/api/players/leaderboards/biggest-builder", playerHandler.GetBiggestBuilder)

	// profiles
	r.Get("/api/profiles/{uuid}", profileHandler.GetProfile)
	r.Get("/api/profiles/{uuid}/daily-mob-kill-chart", profileHandler.GetMobKillChart)
	r.Get("/api/profiles/{uuid}/blocks-broken-pie-chart", profileHandler.GetBlocksBrokenPieChart)
	r.Get("/api/profiles/{uuid}/total-blocks-chart", profileHandler.GetTotalBlocksChart)
	r.Get("/api/profiles/{uuid}/deaths-chart", profileHandler.GetDeathsChart)
	r.Get("/api/profiles/{uuid}/most-killed-mob", profileHandler.GetMostKilledMob)

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
