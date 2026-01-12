package ui

import "fyne.io/fyne/v2/widget"

func SearchBar() *widget.Entry {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Rechercher un artiste...")
	return entry
}
