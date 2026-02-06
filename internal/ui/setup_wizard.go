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

	var continueBtn *widget.Button

	selectBtn := widget.NewButton("Select Folder", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil || list == nil {
				return
			}
			path := list.Path()
			pathLabel.SetText(path)
			pathLabel.TextStyle = fyne.TextStyle{Italic: false}

			valid, msg := system.ValidateGoldenDictPath(path)
			if valid {
				s.ui.State.Config.GoldenDictPath = path
				continueBtn.Enable()
			}
			dialog.ShowInformation("Path Validation", msg, s.ui.Window)
		}, s.ui.Window)
	})

	continueBtn = widget.NewButton("Continue", func() {
		s.ui.ConfigManager.SaveConfig(s.ui.State.Config)
		s.ui.ShowMain()
	})
	continueBtn.Disable()

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		intro,
		container.NewHBox(selectBtn, pathLabel),
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), continueBtn),
	)

	return container.NewPadded(content)
}
