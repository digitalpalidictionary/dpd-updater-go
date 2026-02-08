package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
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
	myApp := app.NewWithID("net.dpdict.dpd-updater")
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
	if !u.ensureGoldenDictClosed() {
		return
	}

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

func (u *UI) ensureGoldenDictClosed() bool {
	gm := system.NewGoldenDictManager()

	running, _ := gm.IsRunning()
	if !running {
		return true
	}

	gm.Close(5 * time.Second)

	running, _ = gm.IsRunning()
	if !running {
		u.showAutoCloseNotification()
		return true
	}

	return u.showGoldenDictBlockingDialog(gm)
}

func (u *UI) showAutoCloseNotification() {
	done := make(chan struct{})
	u.Window.Show()
	u.RunOnMain(func() {
		d := dialog.NewInformation(
			"GoldenDict Closed",
			"GoldenDict was automatically closed to prevent file locks.",
			u.Window,
		)
		d.Show()
		go func() {
			time.Sleep(3 * time.Second)
			d.Hide()
			close(done)
		}()
	})
	<-done
}

func (u *UI) showGoldenDictBlockingDialog(gm *system.GoldenDictManager) bool {
	done := make(chan bool, 1)

	var showDialog func()
	showDialog = func() {
		dialog.ShowConfirm(
			"GoldenDict is Running",
			"GoldenDict is currently running and locks dictionary files.\n\n"+
				"Automatic close failed. Please close GoldenDict manually, then click \"I've closed it\".",
			func(confirmed bool) {
				if confirmed {
					running, _ := gm.IsRunning()
					if !running {
						done <- true
					} else {
						u.RunOnMain(showDialog)
					}
				} else {
					done <- false
				}
			},
			u.Window,
		)
	}

	u.Window.Show()
	u.RunOnMain(showDialog)

	return <-done
}
