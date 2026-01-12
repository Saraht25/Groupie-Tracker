package ui

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func StartApp() {
	a := app.New()
	w := a.NewWindow("Groupie Tracker")
	w.SetContent(Home())
	w.ShowAndRun()
}
