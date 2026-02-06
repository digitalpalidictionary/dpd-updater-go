package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/digitalpalidictionary/dpd-updater-go/internal/config"
)

type UI struct {
	App           fyne.App
	Window        fyne.Window
	State         *AppState
	ConfigManager *config.ConfigManager
}

func NewUI(cfg *config.Config, cm *config.ConfigManager) *UI {
	myApp := app.New()
	myWindow := myApp.NewWindow("DPD Updater")
	myWindow.Resize(fyne.NewSize(800, 600))

	state := NewAppState(cfg)

	return &UI{
		App:           myApp,
		Window:        myWindow,
		State:         state,
		ConfigManager: cm,
	}
}

func (u *UI) Start() {
	if u.State.Config.GoldenDictPath == "" {
		u.ShowSetup()
	} else {
		u.ShowMain()
	}
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
