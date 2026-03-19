package main

import (
	"log"

	"campushub-ctf/internal/config"
	"campushub-ctf/internal/database"
	"campushub-ctf/internal/server"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DBDSN)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	if err := database.EnsureSchema(db); err != nil {
		log.Fatalf("ensure schema: %v", err)
	}

	if err := database.SeedDemoData(db, cfg.UploadDir); err != nil {
		log.Fatalf("seed data: %v", err)
	}

	app := server.New(cfg, db)
	if err := app.Run(); err != nil {
		log.Fatalf("run server: %v", err)
	}
}

