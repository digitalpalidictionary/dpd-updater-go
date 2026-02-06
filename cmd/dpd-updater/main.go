package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("DPD Updater")

	myWindow.SetContent(widget.NewLabel("Hello DPD!"))
	myWindow.ShowAndRun()
}