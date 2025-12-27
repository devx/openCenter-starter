package main

import (
	"log"

	httpadapter "github.com/devx/openCenter-starter/backend/internal/adapters/http"
	"github.com/devx/openCenter-starter/backend/internal/config"
)

func main() {
	cfg := config.Load()
	server, err := httpadapter.New(cfg)
	if err != nil {
		log.Fatalf("server init error: %v", err)
	}

	log.Printf("starting server on %s", cfg.Addr)
	if err := server.Listen(cfg.Addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
