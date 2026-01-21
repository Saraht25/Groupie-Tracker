// Package ui - artist_detail.go affiche la page détaillée d'un artiste.
// Il présente l'image, la biographie, les membres, les lieux de concert avec cartes géographiques,
// et un lien Spotify. C'est la fenêtre d'information complète pour explorer un artiste en profondeur.
package ui

import (
	"Groupie-Tracker/models"
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CreateArtistDetailView construit la page détaillée d'un artiste avec image, membres, lieux et bouton retour
func CreateArtistDetailView(artist models.Artist, app fyne.App, onBack func()) fyne.CanvasObject {
	img := loadDetailImage(artist.Image)

	nameLabel := widget.NewLabel(artist.Name)
	nameLabel.TextStyle.Bold = true
	nameLabel.Alignment = fyne.TextAlignCenter

	creationLabel := widget.NewLabel(
		fmt.Sprintf("Année de création: %d", artist.CreationDate),
	)
	albumLabel := widget.NewLabel(
		fmt.Sprintf("Premier album: %s", artist.FirstAlbum),
	)

	membersLabel := widget.NewLabel("Membres:")
	membersLabel.TextStyle.Bold = true
	membersBox := container.NewVBox()
	for _, member := range artist.Members {
		membersBox.Add(widget.NewLabel("  • " + member))
	}

	locationsLabel := widget.NewLabel("Lieux de concert:")
	locationsLabel.TextStyle.Bold = true
	locationsBox := container.NewVBox()

	for _, loc := range artist.Locations {
		locText := widget.NewLabel("  • " + loc)
		locText.Wrapping = fyne.TextWrapWord

		mapPlaceholder := widget.NewLabel("🗺️ Chargement...")
		mapPlaceholder.Alignment = fyne.TextAlignCenter

		locContainer := container.NewVBox(
			locText,
			mapPlaceholder,
		)

		locationsBox.Add(locContainer)

		go func(loc string, idx int) {
			locMap := createLocationMapForSingle(loc)
			if len(locationsBox.Objects) > idx {
				if vbox, ok := locationsBox.Objects[idx].(*fyne.Container); ok && vbox.Layout != nil {
					if len(vbox.Objects) > 1 {
						vbox.Objects[1] = locMap
						vbox.Refresh()
					}
				}
			}
		}(loc, len(locationsBox.Objects)-1)
	}

	backButton := widget.NewButtonWithIcon("Retour", theme.NavigateBackIcon(), func() {
		onBack()
	})
	backButton.Importance = widget.HighImportance

	spotifyLabel := widget.NewLabel("Écouter sur Spotify:")
	spotifyLabel.TextStyle.Bold = true
	spotifyButton := widget.NewButton("Ouvrir Spotify", func() {
		searchQuery := url.QueryEscape(artist.Name)
		spotifyURL, _ := url.Parse(fmt.Sprintf("https://open.spotify.com/search/%s", searchQuery))
		app.OpenURL(spotifyURL)
	})

	locationsScroll := container.NewVScroll(locationsBox)
	locationsScroll.SetMinSize(fyne.NewSize(600, 400))

	leftCol := container.NewVBox(
		container.NewCenter(img),
		widget.NewSeparator(),
		spotifyLabel,
		container.NewCenter(spotifyButton),
		widget.NewSeparator(),
		locationsLabel,
		locationsScroll,
	)

	membersScroll := container.NewVScroll(membersBox)
	membersScroll.SetMinSize(fyne.NewSize(600, 300))

	rightCol := container.NewVBox(
		nameLabel,
		creationLabel,
		albumLabel,
		widget.NewSeparator(),
		membersLabel,
		membersScroll,
	)

	leftScroll := container.NewVScroll(leftCol)
	rightScroll := container.NewVScroll(rightCol)

	leftScroll.SetMinSize(fyne.NewSize(350, 700))
	rightScroll.SetMinSize(fyne.NewSize(650, 700))

	contentBody := container.NewHBox(
		leftScroll,
		rightScroll,
	)

	header := container.NewHBox(backButton)

	view := container.NewBorder(
		header, nil, nil, nil,
		contentBody,
	)

	return view
}
