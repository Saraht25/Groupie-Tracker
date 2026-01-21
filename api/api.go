// Package api gère la communication avec les APIs externes (GroupieTrackers et Nominatim).
// Il récupère les données des artistes, leurs localisations de concerts et les coordonnées géographiques.
// C'est le bridge essentiel entre l'interface utilisateur et les sources de données externes.
package api

import (
	"Groupie-Tracker/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	baseURL         = "https://groupietrackers.herokuapp.com/api"
	artistsEndpoint = baseURL + "/artists"
)

func fetchAPI(url string, target interface{}) error {
	// Créer client HTTP avec timeout
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	// Faire requête GET
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Vérifier status HTTP
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", resp.Status)
	}

	// Décoder JSON dans target
	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return err
	}
	return nil
}

func GetArtists() ([]models.Artist, error) {
	var artists []models.Artist
	if err := fetchAPI(artistsEndpoint, &artists); err != nil {
		return nil, err
	}

	for i := range artists {
		locs, err := getArtistLocations(artists[i].LocationsURL)
		if err != nil {
			return nil, fmt.Errorf("fetch locations for artist %d: %w", artists[i].Id, err)
		}
		artists[i].Locations = locs
	}

	return artists, nil
}

func getArtistLocations(url string) ([]string, error) {
	var payload struct {
		Locations []string `json:"locations"`
	}
	if err := fetchAPI(url, &payload); err != nil {
		return nil, err
	}
	return payload.Locations, nil
}
