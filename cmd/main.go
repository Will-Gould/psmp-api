package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/Will-Gould/psmp-api/internal/env"
)

func main() {

	godotenv.Load(".env")

	fmt.Println(`
██████╗  █████╗  ██████╗██╗███████╗██╗ ██████╗    ███████╗███╗   ███╗██████╗ 
██╔══██╗██╔══██╗██╔════╝██║██╔════╝██║██╔════╝    ██╔════╝████╗ ████║██╔══██╗
██████╔╝███████║██║     ██║█████╗  ██║██║         ███████╗██╔████╔██║██████╔╝
██╔═══╝ ██╔══██║██║     ██║██╔══╝  ██║██║         ╚════██║██║╚██╔╝██║██╔═══╝ 
██║     ██║  ██║╚██████╗██║██║     ██║╚██████╗    ███████║██║ ╚═╝ ██║██║     
╚═╝     ╚═╝  ╚═╝ ╚═════╝╚═╝╚═╝     ╚═╝ ╚═════╝    ╚══════╝╚═╝     ╚═╝╚═╝     
`)

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "username:password@tcp(localhost:3306)/dummy?parseTime=true"),
		},
	}

	// Structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Database
	db, err := sql.Open("mysql", cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	logger.Info("Connected to database")

	api := application{
		config: cfg,
		db:     db,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("Server failed to start, err: %s", err)
		os.Exit(1)
	}

}
