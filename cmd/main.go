package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
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
		db:   dbConfig{},
	}

	api := application{
		config: cfg,
	}

	// Structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := api.run(api.mount()); err != nil {
		slog.Error("Server failed to start, err: %s", err)
		os.Exit(1)
	}

}
