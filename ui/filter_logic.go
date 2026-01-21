// Package ui - filter_logic.go contient la logique métier de filtrage.
// Il applique les critères de filtrage (dates de création, années d'albums, nombre de membres, localités)
// de manière indépendante de l'UI. C'est un module testable et réutilisable pour tous les filtres.
package ui

import (
	"Groupie-Tracker/models"
	"strconv"
	"strings"
)

// FilterCriteria contient tous les critères de filtrage
type FilterCriteria struct {
	CreationMin    string
	CreationMax    string
	AlbumMin       string
	AlbumMax       string
	MemberCountMin string
	MemberCountMax string
	LocationQuery  string
}

// ApplyFilters applique les critères de filtrage sur la liste d'artistes
func ApplyFilters(artists []models.Artist, criteria FilterCriteria) []models.Artist {
	filtered := artists

	// Filtre année de création
	if criteria.CreationMin != "" || criteria.CreationMax != "" {
		var min, max int
		if criteria.CreationMin != "" {
			min, _ = strconv.Atoi(criteria.CreationMin)
		} else {
			min = 0
		}
		if criteria.CreationMax != "" {
			max, _ = strconv.Atoi(criteria.CreationMax)
		} else {
			max = 9999
		}

		var temp []models.Artist
		for _, a := range filtered {
			if a.CreationDate >= min && a.CreationDate <= max {
				temp = append(temp, a)
			}
		}
		filtered = temp
	}

	// Filtre année du premier album
	if criteria.AlbumMin != "" || criteria.AlbumMax != "" {
		var min, max int
		if criteria.AlbumMin != "" {
			min, _ = strconv.Atoi(criteria.AlbumMin)
		} else {
			min = 0
		}
		if criteria.AlbumMax != "" {
			max, _ = strconv.Atoi(criteria.AlbumMax)
		} else {
			max = 9999
		}

		var temp []models.Artist
		for _, a := range filtered {
			if year, ok := firstYearFromString(a.FirstAlbum); ok {
				if year >= min && year <= max {
					temp = append(temp, a)
				}
			}
		}
		filtered = temp
	}

	// Filtre nombre de membres
	if criteria.MemberCountMin != "" || criteria.MemberCountMax != "" {
		var min, max int
		if criteria.MemberCountMin != "" {
			min, _ = strconv.Atoi(criteria.MemberCountMin)
		} else {
			min = 0
		}
		if criteria.MemberCountMax != "" {
			max, _ = strconv.Atoi(criteria.MemberCountMax)
		} else {
			max = 9999
		}

		var temp []models.Artist
		for _, a := range filtered {
			count := len(a.Members)
			if count >= min && count <= max {
				temp = append(temp, a)
			}
		}
		filtered = temp
	}

	// Filtre par localisation
	locQ := strings.TrimSpace(strings.ToLower(criteria.LocationQuery))
	if locQ != "" {
		var temp []models.Artist
		for _, a := range filtered {
			matched := false
			for _, loc := range a.Locations {
				if strings.Contains(strings.ToLower(loc), locQ) {
					matched = true
					break
				}
			}
			if matched {
				temp = append(temp, a)
			}
		}
		filtered = temp
	}

	return filtered
}
