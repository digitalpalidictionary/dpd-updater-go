package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/digitalpalidictionary/dpd-updater-go/internal/config"
	"github.com/digitalpalidictionary/dpd-updater-go/internal/system"
)

type UI struct {
	App           fyne.App
	Window        fyne.Window
	State         *AppState
	ConfigManager *config.ConfigManager
	Dispatch      chan func()
}

func NewUI(cfg *config.Config, cm *config.ConfigManager) *UI {
	myApp := app.New()
	myApp.SetIcon(resourceIconPng)
	myWindow := myApp.NewWindow("DPD Updater")
	myWindow.Resize(fyne.NewSize(800, 600))

	state := NewAppState(cfg)

	// Add Ctrl-Q (or Cmd-Q on macOS) shortcut to quit
	quitShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyQ, Modifier: fyne.KeyModifierShortcutDefault}
	myWindow.Canvas().AddShortcut(quitShortcut, func(shortcut fyne.Shortcut) {
		myApp.Quit()
	})

	return &UI{
		App:           myApp,
		Window:        myWindow,
		State:         state,
		ConfigManager: cm,
		Dispatch:      make(chan func(), 100),
	}
}

func (u *UI) Start() {
	if u.State.Config.GoldenDictPath != "" {
		// Scan for version
		if version, err := system.ScanForVersion(u.State.Config.GoldenDictPath); err == nil {
			if version != u.State.Config.InstalledVersion {
				u.State.Config.InstalledVersion = version
				u.ConfigManager.SaveConfig(u.State.Config)
			}
		}
	}

	if u.State.Config.GoldenDictPath == "" {
		u.ShowSetup()
	} else {
		u.ShowMain()
	}

	// Start dispatcher loop in a goroutine
	go func() {
		for f := range u.Dispatch {
			f()
		}
	}()

	u.Window.ShowAndRun()
}

func (u *UI) ShowSetup() {
	setup := NewSetupWizard(u)
	u.Window.SetContent(setup.Render())
}

func (u *UI) ShowMain() {
	main := NewMainWindow(u)
	u.Window.SetContent(main.Render())
}

func (u *UI) RunOnMain(f func()) {
	u.Dispatch <- f
}