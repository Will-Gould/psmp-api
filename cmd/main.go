package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"

	"github.com/Will-Gould/psmp-api/assets"
	"github.com/Will-Gould/psmp-api/internal/env"
)

func main() {

	godotenv.Load(".env")

	assets.PrintLogo()

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "username:password@tcp(localhost:3306)/dummy?parseTime=true"),
		},
	}

	// Structured logging
	w := os.Stderr
	// logger := slog.New(tint.NewHandler(w, nil))
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		},
		),
	))

	// Database
	db, err := sql.Open("mysql", cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	slog.Log(context.Background(), slog.LevelInfo, "Connected to database")

	api := application{
		config: cfg,
		db:     db,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("Server failed to start, err: %s", err)
		os.Exit(1)
	}

}
