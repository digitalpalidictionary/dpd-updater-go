package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/digitalpalidictionary/dpd-updater-go/internal/github"
	"github.com/digitalpalidictionary/dpd-updater-go/internal/installer"
	"github.com/digitalpalidictionary/dpd-updater-go/internal/system"
)

type MainWindow struct {
	ui *UI
}

func NewMainWindow(ui *UI) *MainWindow {
	return &MainWindow{ui: ui}
}

func (m *MainWindow) Render() fyne.CanvasObject {
	u := m.ui
	
	title := widget.NewLabelWithStyle("DPD Updater", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	currentVersion := widget.NewLabel(fmt.Sprintf("Installed: %s", u.State.Config.InstalledVersion))
	latestVersion := widget.NewLabel("Latest: Checking...")
	statusLabel := widget.NewLabel("")

	logList := widget.NewList(
		func() int { return len(u.State.Logs) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(u.State.Logs[id])
		},
	)

	progress := widget.NewProgressBar()
	progress.Hide()

	var updateBtn *widget.Button
	var checkBtn *widget.Button

	checkUpdates := func() {
		u.State.SetProcessing(true)
		checkBtn.Disable()
		statusLabel.SetText("Checking for updates...")
		
		go func() {
			client := github.NewGitHubClient()
			release, err := client.GetLatestRelease()
			
			u.State.Lock()
			if err != nil {
				u.State.Logs = append(u.State.Logs, fmt.Sprintf("Error checking updates: %v", err))
				u.State.Unlock()
				u.Window.Canvas().Refresh(statusLabel)
				return
			}
			u.State.LatestRelease = release
			u.State.Unlock()

			latestVersion.SetText(fmt.Sprintf("Latest: %s", release.Version))
			
			comp := github.CompareVersions(u.State.Config.InstalledVersion, release.Version)
			if comp < 0 {
				statusLabel.SetText("Update available!")
				updateBtn.Enable()
			} else {
				statusLabel.SetText("You are up to date.")
				updateBtn.Disable()
			}
			
			u.State.SetProcessing(false)
			checkBtn.Enable()
		}()
	}

	checkBtn = widget.NewButton("Check for Updates", checkUpdates)

	updateBtn = widget.NewButton("Update Now", func() {
		u.State.SetProcessing(true)
		updateBtn.Disable()
		checkBtn.Disable()
		progress.Show()

		go func() {
			gm := system.NewGoldenDictManager()
			if running, _ := gm.IsRunning(); running {
				u.State.AddLog("Closing GoldenDict...")
				gm.Close(10 * time.Second)
			}

			inst := installer.NewInstaller(u.State.Config, func(msg string, p int) {
				u.State.SetStatus(msg, float64(p)/100.0)
				progress.SetValue(float64(p)/100.0)
				u.Window.Canvas().Refresh(logList)
				logList.ScrollToBottom()
			})

			// 1. Create temp dir for download
			tempDir := filepath.Join(u.State.Config.GoldenDictPath, "_dpd_download_temp")
			os.MkdirAll(tempDir, 0755)
			defer os.RemoveAll(tempDir)

			// 2. Download
			zipPath, err := inst.DownloadRelease(u.State.LatestRelease.AssetURL, tempDir)
			if err != nil {
				u.State.AddLog(fmt.Sprintf("Download failed: %v", err))
				u.State.SetProcessing(false)
				checkBtn.Enable()
				return
			}

			// 3. Backup
			inst.BackupExisting(u.State.Config.GoldenDictPath)

			// 4. Install
			err = inst.InstallUpdate(zipPath, u.State.Config.GoldenDictPath)
			if err != nil {
				u.State.AddLog(fmt.Sprintf("Installation failed: %v", err))
			} else {
				u.State.Config.InstalledVersion = u.State.LatestRelease.Version
				u.ConfigManager.SaveConfig(u.State.Config)
				currentVersion.SetText(fmt.Sprintf("Installed: %s", u.State.Config.InstalledVersion))
				u.State.AddLog("Update complete!")
				
				u.State.AddLog("Restarting GoldenDict...")
				gm.Reopen()
			}

			u.State.SetProcessing(false)
			checkBtn.Enable()
			progress.Hide()
		}()
	})
	updateBtn.Disable()

	top := container.NewVBox(
		title,
		widget.NewSeparator(),
		container.NewVBox(
			widget.NewLabelWithStyle("Version Information", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			currentVersion,
			latestVersion,
			statusLabel,
		),
	)

	buttons := container.NewHBox(layout.NewSpacer(), checkBtn, updateBtn, layout.NewSpacer())

	content := container.NewBorder(
		top,
		container.NewVBox(progress, buttons),
		nil,
		nil,
		logList,
	)

	// Initial check
	if u.State.Config.AutoCheckUpdates {
		go checkUpdates()
	}

	return container.NewPadded(content)
}
