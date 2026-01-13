package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"groupie-tracker/api"
)

func StartApp() {
	a := app.New()
	w := a.NewWindow("Groupie Tracker")

	artists, err := api.GetArtists()
	if err != nil {
		w.SetContent(widget.NewLabel("Erreur chargement API"))
		w.ShowAndRun()
		return
	}

	var cards []fyne.CanvasObject
	for _, artist := range artists {
		card := widget.NewCard(
			artist.Name,
			artist.FirstAlbum,
			widget.NewLabel(
				"Members: " + 
				fmt.Sprint(artist.Members),
			),
		)
		cards = append(cards, card)
	}

	w.SetContent(container.NewVScroll(
		container.NewVBox(cards...),
	))

	w.Resize(fyne.NewSize(1200, 800))
	w.ShowAndRun()
}
