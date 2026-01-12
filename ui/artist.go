package ui

import (
	"fyne.io/fyne/v2/widget"
	"Groupie-Tracker/models"
)

func ArtistCard(a models.Artist) *widget.Card {
	return widget.NewCard(
		a.Name,
		"",
		widget.NewLabel("Membres: "+a.Members),
	)
}
