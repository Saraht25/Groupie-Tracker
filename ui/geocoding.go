// Package ui - geocoding.go gère la conversion des noms de localités en coordonnées géographiques.
// Il utilise l'API Nominatim (OpenStreetMap) pour transformer "Paris-France" en lat/lon.
// C'est essentiel pour afficher les lieux de concerts sur les cartes à la bonne position.
package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// geocodeLocation interroge Nominatim pour convertir une adresse en lat/lon
func geocodeLocation(query string) (lat string, lon string, ok bool) {
	norm := normalizeLocationQuery(query)
	endpoint := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1", url.QueryEscape(norm))
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", "", false
	}
	req.Header.Set("User-Agent", "GroupieTracker/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", false
	}

	var results []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return "", "", false
	}
	if len(results) == 0 {
		return "", "", false
	}
	return results[0].Lat, results[0].Lon, true
}

// normalizeLocationQuery nettoie une localisation "ville-pays" ou avec underscores en "ville, pays"
func normalizeLocationQuery(loc string) string {
	parts := strings.FieldsFunc(loc, func(r rune) bool {
		return r == '-' || r == '_'
	})
	if len(parts) == 0 {
		return strings.TrimSpace(loc)
	}
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return strings.Join(parts, ", ")
}
