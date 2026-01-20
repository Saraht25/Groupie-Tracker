package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func StartApp() {
	a := app.New()
	w := a.NewWindow("Groupie Tracker")

	// Définir une taille de fenêtre appropriée pour une application de bureau
	w.Resize(fyne.NewSize(1200, 800))
	w.CenterOnScreen()

	w.SetContent(Home(a, w))
	w.ShowAndRun()
}
