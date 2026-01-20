package ui

import (
	"Groupie-Tracker/models"
	"errors"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func ArtistCard(a models.Artist) *widget.Card {
	return widget.NewCard(
		a.Name,
		"",
		widget.NewLabel("Membres: "+strings.Join(a.Members, ", ")),
	)
}

func OpenSpotifyURL(rawURL string) error {
	if strings.TrimSpace(rawURL) == "" {
		return errors.New("spotify URL vide")
	}

	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return errors.New("spotify URL invalide")
	}

	app := fyne.CurrentApp()
	if app == nil {
		return errors.New("application Fyne non initialisée")
	}
	return app.OpenURL(parsed)
}
