package main

import (
	"log"

	"github.com/digitalpalidictionary/dpd-updater-go/internal/config"
	"github.com/digitalpalidictionary/dpd-updater-go/internal/ui"
)

func main() {
	cm, err := config.NewConfigManager()
	if err != nil {
		log.Fatalf("Failed to create config manager: %v", err)
	}

	cfg, err := cm.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	application := ui.NewUI(cfg, cm)
	application.Start()
}
