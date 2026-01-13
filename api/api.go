package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"groupie-tracker/models"
)

const (
	baseURL           = "https://groupietrackers.herokuapp.com/api"
	artistsEndpoint   = baseURL + "/artists"
	locationsEndpoint = baseURL + "/locations"
	datesEndpoint     = baseURL + "/dates"
	relationEndpoint  = baseURL + "/relation"
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
	err := fetchAPI(artistsEndpoint, &artists)
	return artists, err
}

func GetRelations() ([]models.Relation, error) {
	var relations []models.Relation
	err := fetchAPI(relationEndpoint, &relations)
	return relations, err
}

func GetLocations() ([]models.Location, error) {
	var locations []models.Location
	err := fetchAPI(locationsEndpoint, &locations)
	return locations, err
}

func GetDates() ([]models.Date, error) {
	var dates []models.Date
	err := fetchAPI(datesEndpoint, &dates)
	return dates, err
}
