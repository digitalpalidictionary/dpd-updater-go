package ui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/digitalpalidictionary/dpd-updater-go/internal/github"
	"github.com/digitalpalidictionary/dpd-updater-go/internal/installer"
	"github.com/digitalpalidictionary/dpd-updater-go/internal/system"
)

type MainWindow struct {
	ui           *UI
	updateCancel context.CancelFunc
}

func NewMainWindow(ui *UI) *MainWindow {
	return &MainWindow{ui: ui}
}

func (m *MainWindow) runOnMain(f func()) {
	fyne.Do(f)
}

func (m *MainWindow) Render() fyne.CanvasObject {
	u := m.ui

	title := widget.NewLabelWithStyle("DPD Updater", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	currentVerBind := binding.NewString()
	currentVerBind.Set(fmt.Sprintf("Installed: %s", u.State.Config.InstalledVersion))
	
	latestVerBind := binding.NewString()
	latestVerBind.Set("Latest: Checking...")
	
	statusBind := binding.NewString()
	statusBind.Set("")

	progressBind := binding.NewFloat()
	progressBind.Set(0.0)

	currentVersion := widget.NewLabelWithData(currentVerBind)
	latestVersion := widget.NewLabelWithData(latestVerBind)
	statusLabel := widget.NewLabelWithData(statusBind)

	logList := widget.NewList(
		func() int { return len(u.State.Logs) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			u.State.RLock()
			defer u.State.RUnlock()
			if id < len(u.State.Logs) {
				obj.(*widget.Label).SetText(u.State.Logs[id])
			}
		},
	)

	progress := widget.NewProgressBarWithData(progressBind)
	progress.Hide()

	var updateBtn *widget.Button
	var checkBtn *widget.Button
	var cancelBtn *widget.Button

	cancelBtn = widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		if m.updateCancel != nil {
			m.updateCancel()
			u.State.AddLog("Update cancelled by user.")
			m.runOnMain(func() {
				logList.Refresh()
			})
		}
	})
	cancelBtn.Hide()

	checkUpdates := func() {
		u.State.SetProcessing(true)
		m.runOnMain(func() { checkBtn.Disable() })
		statusBind.Set("Checking for updates...")

		go func() {
			client := github.NewGitHubClient()
			release, err := client.GetLatestRelease()

			if err != nil {
				u.State.AddLog(fmt.Sprintf("Error checking updates: %v", err))
				statusBind.Set("Error checking updates")
				u.State.SetProcessing(false)
				m.runOnMain(func() {
					checkBtn.Enable()
					logList.Refresh()
				})
				return
			}

			u.State.Lock()
			u.State.LatestRelease = release
			u.State.Unlock()

			latestVerBind.Set(fmt.Sprintf("Latest: %s", release.Version))

			comp := github.CompareVersions(u.State.Config.InstalledVersion, release.Version)
			
			u.State.SetProcessing(false)
			m.runOnMain(func() {
				if comp < 0 {
					statusBind.Set("Update available!")
					updateBtn.Enable()
				} else {
					statusBind.Set("You are up to date.")
					updateBtn.Disable()
				}
				checkBtn.Enable()
				logList.Refresh()
			})
		}()
	}

	checkBtn = widget.NewButton("Check for Updates", checkUpdates)

	updateBtn = widget.NewButton("Update Now", func() {
		u.State.SetProcessing(true)
		
		ctx, cancel := context.WithCancel(context.Background())
		m.updateCancel = cancel

		m.runOnMain(func() {
			updateBtn.Disable()
			checkBtn.Disable()
			progress.Show()
			cancelBtn.Show()
		})

		go func() {
			defer func() {
				m.runOnMain(func() {
					cancelBtn.Hide()
					progress.Hide()
					checkBtn.Enable()
				})
				u.State.SetProcessing(false)
			}()

			gm := system.NewGoldenDictManager()
			if running, _ := gm.IsRunning(); running {
				u.State.AddLog("Closing GoldenDict...")
				m.runOnMain(func() { logList.Refresh() })
				gm.Close(10 * time.Second)
			}

			inst := installer.NewInstaller(u.State.Config, func(msg string, p int) {
				u.State.SetStatus(msg, float64(p)/100.0)
				m.runOnMain(func() {
					progressBind.Set(float64(p) / 100.0)
					logList.Refresh()
					logList.ScrollToBottom()
				})
			})

			tempDir := filepath.Join(u.State.Config.GoldenDictPath, "_dpd_download_temp")
			os.MkdirAll(tempDir, 0755)
			defer os.RemoveAll(tempDir)

			// 2. Download
			zipPath, err := inst.DownloadRelease(ctx, u.State.LatestRelease.AssetURL, tempDir)
			if err != nil {
				if err == context.Canceled {
					u.State.AddLog("Download cancelled.")
				} else {
					u.State.AddLog(fmt.Sprintf("Download failed: %v", err))
				}
				m.runOnMain(func() { logList.Refresh() })
				return
			}

			// 3. Backup
			_, err = inst.BackupExisting(ctx, u.State.Config.GoldenDictPath)
			if err != nil {
				if err == context.Canceled {
					u.State.AddLog("Backup cancelled.")
				} else {
					u.State.AddLog(fmt.Sprintf("Backup failed: %v", err))
				}
				m.runOnMain(func() { logList.Refresh() })
				return
			}

			// 4. Install
			err = inst.InstallUpdate(ctx, zipPath, u.State.Config.GoldenDictPath)
			if err != nil {
				if err == context.Canceled {
					u.State.AddLog("Installation cancelled.")
				} else {
					u.State.AddLog(fmt.Sprintf("Installation failed: %v", err))
				}
				m.runOnMain(func() { logList.Refresh() })
			} else {
				u.State.Lock()
				u.State.Config.InstalledVersion = u.State.LatestRelease.Version
				u.ConfigManager.SaveConfig(u.State.Config)
				u.State.Unlock()
				
				currentVerBind.Set(fmt.Sprintf("Installed: %s", u.State.Config.InstalledVersion))
				u.State.AddLog("Update complete!")
				u.State.AddLog("Restarting GoldenDict...")
				gm.Reopen()
				m.runOnMain(func() { logList.Refresh() })
			}
		}()
	})
	updateBtn.Disable()

	settingsBtn := widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), func() {
		m.showSettings()
	})

	top := container.NewVBox(
		container.NewHBox(title, layout.NewSpacer(), settingsBtn),
		widget.NewSeparator(),
		container.NewVBox(
			widget.NewLabelWithStyle("Version Information", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			currentVersion,
			latestVersion,
			statusLabel,
		),
	)

	buttons := container.NewHBox(layout.NewSpacer(), checkBtn, updateBtn, layout.NewSpacer())

	// Progress bar and Cancel button row
	progressRow := container.NewBorder(nil, nil, nil, cancelBtn, progress)

	content := container.NewBorder(
		top,
		container.NewVBox(progressRow, buttons),
		nil,
		nil,
		logList,
	)

	if u.State.Config.AutoCheckUpdates {
		go checkUpdates()
	}

	return container.NewPadded(content)
}

func (m *MainWindow) showSettings() {
	u := m.ui

	pathLabel := widget.NewLabel(u.State.Config.GoldenDictPath)
	pathLabel.Wrapping = fyne.TextWrapBreak

	changePathBtn := widget.NewButton("Change Folder", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil || list == nil {
				return
			}
			path := list.Path()
			valid, msg := system.ValidateGoldenDictPath(path)
			if valid {
				u.State.Config.GoldenDictPath = path
				m.runOnMain(func() {
					pathLabel.SetText(path)
				})
				u.ConfigManager.SaveConfig(u.State.Config)
			} else {
				dialog.ShowError(fmt.Errorf(msg), u.Window)
			}
		}, u.Window)
	})

	autoCheck := widget.NewCheck("Check for updates on startup", func(val bool) {
		u.State.Config.AutoCheckUpdates = val
		u.ConfigManager.SaveConfig(u.State.Config)
	})
	autoCheck.SetChecked(u.State.Config.AutoCheckUpdates)

	backup := widget.NewCheck("Create backup before updating", func(val bool) {
		u.State.Config.BackupBeforeUpdate = val
		u.ConfigManager.SaveConfig(u.State.Config)
	})
	backup.SetChecked(u.State.Config.BackupBeforeUpdate)

	content := container.NewVBox(
		widget.NewLabelWithStyle("GoldenDict Folder:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		pathLabel,
		changePathBtn,
		widget.NewSeparator(),
		autoCheck,
		backup,
	)

	dialog.ShowCustom("Settings", "Close", content, u.Window)
}
