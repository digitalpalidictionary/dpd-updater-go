package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/digitalpalidictionary/dpd-updater-go/internal/system"
)

type SetupWizard struct {
	ui *UI
}

func NewSetupWizard(ui *UI) *SetupWizard {
	return &SetupWizard{ui: ui}
}

func (s *SetupWizard) Render() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Welcome to DPD Updater", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	// title.TextSize = 24 removed as it's not a field on Label

	intro := widget.NewLabel("This wizard will help you set up the DPD Updater for your GoldenDict installation.\n\nPlease select your GoldenDict content/dictionaries folder.")
	intro.Wrapping = fyne.TextWrapWord

	pathLabel := widget.NewLabel("No folder selected")
	pathLabel.TextStyle = fyne.TextStyle{Italic: true}

	statusLabel := widget.NewLabel("")
	statusLabel.Wrapping = fyne.TextWrapWord

	var continueBtn *widget.Button

	selectBtn := widget.NewButton("Select Folder", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil || list == nil {
				return
			}
			path := list.Path()
			pathLabel.SetText(path)
			pathLabel.TextStyle = fyne.TextStyle{Italic: false}
			statusLabel.SetText("") // Clear auto-detect status on manual selection

			valid, msg := system.ValidateGoldenDictPath(path)
			if valid {
				s.ui.State.Config.GoldenDictPath = path
				continueBtn.Enable()
			} else {
				continueBtn.Disable()
			}
			dialog.ShowInformation("Path Validation", msg, s.ui.Window)
		}, s.ui.Window)
	})

	continueBtn = widget.NewButton("Continue", func() {
		s.ui.ConfigManager.SaveConfig(s.ui.State.Config)
		s.ui.ShowMain()
	})
	continueBtn.Disable()

	// Auto-detection logic
	go func() {
		configPath, err := system.GetGoldenDictConfigPath()
		if err != nil {
			statusLabel.SetText("Could not find GoldenDict config.")
			return
		}

		paths, err := system.ParseGoldenDictPaths(configPath)
		if err != nil {
			statusLabel.SetText("Could not read GoldenDict config.")
			return
		}

		suggested := system.AnalyzeGoldenDictPaths(paths)

		if suggested != "" {
			pathLabel.SetText(suggested)
			pathLabel.TextStyle = fyne.TextStyle{Italic: false}

			valid, _ := system.ValidateGoldenDictPath(suggested)

			if valid {
				s.ui.State.Config.GoldenDictPath = suggested
				continueBtn.Enable()
				statusLabel.SetText("✓ Auto-detected from GoldenDict settings.")
			} else {
				// If validation fails, we don't pre-fill or enable, but we don't show an error either
				// to avoid confusing the user with a broken path.
				pathLabel.SetText("No folder selected")
				pathLabel.TextStyle = fyne.TextStyle{Italic: true}
			}
		} else {
			if len(paths) > 0 {
				statusLabel.SetText("ℹ Tip: Organize your dictionaries into a single master folder (e.g., Documents/GoldenDict) for easier management.")
			}
		}
	}()

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		intro,
		container.NewHBox(selectBtn, pathLabel),
		statusLabel,
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), continueBtn),
	)

	return container.NewPadded(content)
}
