package ui

import (
	"context"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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

	latestVerBind := binding.NewString()
	latestVerBind.Set("Latest: Checking...")

	statusBind := binding.NewString()
	statusBind.Set("Ready")

	progressBind := binding.NewFloat()
	progressBind.Set(0.0)

	progressTextBind := binding.NewString()
	progressTextBind.Set("")

	// Dynamic Version Container

	versionInfoContainer := container.New(layout.NewCustomPaddedVBoxLayout(0))

	latestVersion := widget.NewLabelWithData(latestVerBind)

	latestVersion.Selectable = true

	statusLabel := widget.NewLabelWithData(statusBind)

	statusLabel.Selectable = true

	statusLabel.TextStyle = fyne.TextStyle{Italic: true}

	progress := widget.NewProgressBarWithData(progressBind)

	progressLabel := widget.NewLabelWithData(progressTextBind)

	progressLabel.Selectable = true

	// Horizontal progress row: [ ProgressBar ] [ MB / MB ]

	progressRow := container.NewBorder(nil, nil, nil, progressLabel, progress)

	progressRow.Hide()

	var updateBtn *widget.Button
	var checkBtn *widget.Button
	var cancelBtn *widget.Button
	var cleanupBtn *widget.Button

	// Helper to create selectable label-like widgets
	newSelectableLabel := func(text string, italic bool) *widget.Label {
		l := widget.NewLabel(text)
		l.Selectable = true
		if italic {
			l.TextStyle = fyne.TextStyle{Italic: true}
		}
		return l
	}

	refreshVersionUI := func() {
		versionInfoContainer.Objects = nil

		if len(u.State.DPDInstances) == 0 {
			versionInfoContainer.Add(newSelectableLabel("No DPD dictionaries found.", false))
			versionInfoContainer.Add(latestVersion)
			versionInfoContainer.Add(statusLabel)
			versionInfoContainer.Refresh()
			return
		}

		// 1. Identify which filenames are actually duplicated
		counts := make(map[string]int)
		for _, inst := range u.State.DPDInstances {
			counts[filepath.Base(inst.Path)]++
		}

		// 2. Filter to only include instances of duplicated filenames
		var duplicates []system.DPDInfo
		for _, inst := range u.State.DPDInstances {
			if counts[filepath.Base(inst.Path)] > 1 {
				duplicates = append(duplicates, inst)
			}
		}

		if len(duplicates) > 0 {
			header := canvas.NewText("Duplicate Installations Detected!", color.RGBA{R: 200, G: 0, B: 0, A: 255})
			header.TextStyle = fyne.TextStyle{Bold: true}
			versionInfoContainer.Add(header)

			// 3. Group duplicates by date
			dateGroups := make(map[string][]system.DPDInfo)
			var dates []string
			for _, inst := range duplicates {
				d := inst.Date.Format("2006-01-02")
				if _, exists := dateGroups[d]; !exists {
					dates = append(dates, d)
					dateGroups[d] = []system.DPDInfo{}
				}
				dateGroups[d] = append(dateGroups[d], inst)
			}

			sort.Slice(dates, func(i, j int) bool { return dates[i] > dates[j] })

			for i, d := range dates {
				status := " [Old/Duplicate]"
				if i == 0 {
					status = " [Active]"
				}

				dateHeader := widget.NewLabelWithStyle(d+status, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
				versionInfoContainer.Add(dateHeader)

				// Group by directory within this date to reduce lines
				pathGroups := make(map[string][]string)
				var paths []string
				for _, inst := range dateGroups[d] {
					dir := filepath.Dir(inst.Path)
					if _, exists := pathGroups[dir]; !exists {
						paths = append(paths, dir)
					}
					pathGroups[dir] = append(pathGroups[dir], filepath.Base(inst.Path))
				}

				for _, path := range paths {
					txt := fmt.Sprintf("  • %s (%s)", path, strings.Join(pathGroups[path], ", "))
					versionInfoContainer.Add(newSelectableLabel(txt, false))
				}
			}

			if cleanupBtn != nil {
				cleanupBtn.Show()
			}
		} else {
			// No duplicates, show simple installed version
			versionInfoContainer.Add(newSelectableLabel(fmt.Sprintf("Installed: %s", u.State.Config.InstalledVersion), false))
			if cleanupBtn != nil {
				cleanupBtn.Hide()
			}
		}

		versionInfoContainer.Add(latestVersion)
		versionInfoContainer.Add(statusLabel)
		versionInfoContainer.Refresh()
	}

	performDuplicateCheck := func() {
		if u.State.Config.GoldenDictPath != "" {
			if instances, err := system.FindAllDPDInstances(u.State.Config.GoldenDictPath); err == nil {
				u.State.Lock()
				u.State.DPDInstances = instances

				// Find newest to set as "InstalledVersion"
				var newest system.DPDInfo
				for _, inst := range instances {
					if newest.Date.IsZero() || inst.Date.After(newest.Date) {
						newest = inst
					}
				}

				if !newest.Date.IsZero() {
					ver := newest.Date.Format("2006-01-02")
					u.State.Config.InstalledVersion = ver
					u.ConfigManager.SaveConfig(u.State.Config)
				}
				u.State.Unlock()
				m.runOnMain(refreshVersionUI)
			}
		}
	}

	cancelBtn = widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		if m.updateCancel != nil {
			m.updateCancel()
		}
	})
	cancelBtn.Hide()

	cleanupBtn = widget.NewButton("Delete Old Versions", func() {
		u.State.Lock()
		instances := u.State.DPDInstances
		u.State.Unlock()

		if len(instances) <= 1 {
			return
		}

		groups := make(map[string][]system.DPDInfo)
		for _, inst := range instances {
			filename := filepath.Base(inst.Path)
			groups[filename] = append(groups[filename], inst)
		}

		var toDelete []string
		var keepDetails []string

		for _, group := range groups {
			if len(group) <= 1 {
				continue
			}

			var keep system.DPDInfo
			for _, inst := range group {
				if keep.Date.IsZero() || inst.Date.After(keep.Date) {
					keep = inst
				}
			}

			for _, inst := range group {
				if inst.Path != keep.Path {
					toDelete = append(toDelete, inst.Path)
				}
			}
			keepDetails = append(keepDetails, fmt.Sprintf("%s (%s) at %s", keep.Bookname, keep.Date.Format("2006-01-02"), filepath.Dir(keep.Path)))
		}

		if len(toDelete) == 0 {
			return
		}

		msg := fmt.Sprintf("Delete %d old dictionary versions?\n\nKeeping newest versions for:\n- %s",
			len(toDelete), strings.Join(keepDetails, "\n- "))

		dialog.ShowConfirm("Confirm Deletion", msg, func(ok bool) {
			if ok {
				statusBind.Set("Deleting old versions...")
				err := installer.DeleteFolders(toDelete)
				if err != nil {
					dialog.ShowError(err, u.Window)
					statusBind.Set("Deletion failed.")
				} else {
					statusBind.Set("Deletion complete.")
					// Trigger check updates to refresh
					checkBtn.OnTapped()
				}
			}
		}, u.Window)
	})
	cleanupBtn.Hide()

	checkUpdates := func() {
		u.State.SetProcessing(true)
		m.runOnMain(func() { checkBtn.Disable() })
		statusBind.Set("Checking for updates...")

		go func() {
			// Check and close GoldenDict before checking for updates
			gm := system.NewGoldenDictManager()
			if running, _ := gm.IsRunning(); running {
				statusBind.Set("Closing GoldenDict...")
				gm.Close(5 * time.Second)

				// Verify it closed
				if running, _ = gm.IsRunning(); !running {
					u.showAutoCloseNotification()
				}
			}

			performDuplicateCheck()

			client := github.NewGitHubClient()
			release, err := client.GetLatestRelease()

			if err != nil {
				statusBind.Set(fmt.Sprintf("Error: %v", err))
				u.State.SetProcessing(false)
				m.runOnMain(func() { checkBtn.Enable() })
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
			progressRow.Show()
			cancelBtn.Show()
		})

		go func() {
			defer func() {
				m.runOnMain(func() {
					cancelBtn.Hide()
					progressRow.Hide()
					checkBtn.Enable()
				})
				u.State.SetProcessing(false)
			}()

			gm := system.NewGoldenDictManager()
			if running, _ := gm.IsRunning(); running {
				statusBind.Set("Closing GoldenDict...")
				gm.Close(10 * time.Second)

				// Verify it closed and show notification
				if running, _ = gm.IsRunning(); !running {
					u.showAutoCloseNotification()
				}
			}

			inst := installer.NewInstaller(u.State.Config, func(msg string, p int) {
				m.runOnMain(func() {
					progressBind.Set(float64(p) / 100.0)

					if strings.HasPrefix(msg, "Downloading...") {
						// Extract just the numbers: "5.1 / 257.7 MB"
						displayMsg := strings.TrimPrefix(msg, "Downloading...")
						progressTextBind.Set(strings.TrimSpace(displayMsg))
						statusBind.Set("Downloading update...")
					} else {
						statusBind.Set(msg)
						progressTextBind.Set("")
					}
				})
			})

			tempDir := filepath.Join(u.State.Config.GoldenDictPath, "_dpd_download_temp")
			os.MkdirAll(tempDir, 0755)
			defer os.RemoveAll(tempDir)

			zipPath, err := inst.DownloadRelease(ctx, u.State.LatestRelease.AssetURL, tempDir)
			if err != nil {
				if err == context.Canceled {
					statusBind.Set("Update cancelled.")
				} else {
					statusBind.Set(fmt.Sprintf("Download failed: %v", err))
				}
				return
			}

			_, err = inst.BackupExisting(ctx, u.State.Config.GoldenDictPath)
			if err != nil {
				if err == context.Canceled {
					statusBind.Set("Update cancelled.")
				} else {
					statusBind.Set(fmt.Sprintf("Backup failed: %v", err))
				}
				return
			}

			err = inst.InstallUpdate(ctx, zipPath, u.State.Config.GoldenDictPath)
			if err != nil {
				if err == context.Canceled {
					statusBind.Set("Update cancelled.")
				} else {
					statusBind.Set(fmt.Sprintf("Installation failed: %v", err))
				}
			} else {
				u.State.Lock()
				u.State.Config.InstalledVersion = u.State.LatestRelease.Version
				u.ConfigManager.SaveConfig(u.State.Config)
				u.State.Unlock()

				m.runOnMain(performDuplicateCheck)
				statusBind.Set("Update complete! Restarting GoldenDict...")
				gm.Reopen()
			}
		}()
	})
	updateBtn.Disable()

	settingsBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		m.showSettings()
	})

	// Initial refresh
	refreshVersionUI()

	top := container.NewVBox(
		container.NewHBox(title, layout.NewSpacer(), settingsBtn),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Version Information", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		versionInfoContainer,
	)

	buttons := container.NewHBox(layout.NewSpacer(), checkBtn, updateBtn, cancelBtn, cleanupBtn, layout.NewSpacer())

	// Clean layout: Version info at top, progress/buttons at bottom
	content := container.NewBorder(
		top,
		container.NewVBox(progressRow, buttons),
		nil,
		nil,
		layout.NewSpacer(), // Empty center
	)

	if u.State.Config.AutoCheckUpdates {
		go checkUpdates()
	} else {
		go performDuplicateCheck()
	}

	return container.NewPadded(content)
}

func (m *MainWindow) showSettings() {
	u := m.ui

	pathLabel := widget.NewLabel(u.State.Config.GoldenDictPath)
	pathLabel.Wrapping = fyne.TextWrapBreak

	changePathBtn := widget.NewButton("Change Folder", func() {
		// Ensure GoldenDict is closed before opening file dialog to avoid file locks
		gm := system.NewGoldenDictManager()
		if running, _ := gm.IsRunning(); running {
			gm.Close(5 * time.Second)
		}

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
