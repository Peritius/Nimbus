package main

import (
	"log"

	"github.com/Peritius/Nimbus/internal/config"
	"github.com/Peritius/Nimbus/internal/tray"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	tray.Run(cfg)
}
