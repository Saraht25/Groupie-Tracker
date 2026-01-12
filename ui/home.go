package ui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Home() *container.Container {
	title := widget.NewLabel("Bienvenue sur Groupie Tracker")
	return container.NewVBox(title)
}
