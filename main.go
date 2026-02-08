package main

import (
	"github.com/digitalpalidictionary/dpd-updater-go/internal/config"
	"github.com/digitalpalidictionary/dpd-updater-go/internal/ui"
)

func main() {
	cm := config.NewConfigManager()
	cfg := cm.LoadConfig()

	application := ui.NewUI(cfg, cm)
	application.Start()
}
