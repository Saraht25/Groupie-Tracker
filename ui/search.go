package ui

import (
	"regexp"
	"strconv"
	"strings"
	"Groupie-Tracker/models"
)

type SuggestionType string

const (
	SuggestionArtist      SuggestionType = "artist"
	SuggestionMember      SuggestionType = "member"
	SuggestionLocation    SuggestionType = "location"
	SuggestionFirstAlbum  SuggestionType = "first_album_year"
	SuggestionCreationYear SuggestionType = "creation_year"
)

type Suggestion struct {
	Label    string
	Type     SuggestionType
	ArtistID int
}

func SearchArtists(query string, artists []models.Artist) []models.Artist {
	q := strings.TrimSpace(strings.ToLower(query))
	if q == "" {
		return artists
	}

	var result []models.Artist
	for _, a := range artists {
		if matchArtist(q, a) {
			result = append(result, a)
		}
	}
	return result
}

func BuildSuggestions(query string, artists []models.Artist) []Suggestion {
	q := strings.TrimSpace(strings.ToLower(query))
	if q == "" {
		return nil
	}

	seen := make(map[string]struct{})
	var suggestions []Suggestion

	add := func(label string, t SuggestionType, artistID int) {
		key := string(t) + "|" + strings.ToLower(label)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		suggestions = append(suggestions, Suggestion{
			Label:    label,
			Type:     t,
			ArtistID: artistID,
		})
	}

	for _, a := range artists {
		if containsIgnoreCase(a.Name, q) {
			add(a.Name, SuggestionArtist, a.ID)
		}
		for _, m := range a.Members {
			if containsIgnoreCase(m, q) {
				add(m, SuggestionMember, a.ID)
			}
		}
		for _, loc := range a.Locations {
			if containsIgnoreCase(loc, q) {
				add(loc, SuggestionLocation, a.ID)
			}
		}
		if year, ok := firstYearFromString(a.FirstAlbum); ok && strings.Contains(strconv.Itoa(year), q) {
			add(strconv.Itoa(year), SuggestionFirstAlbum, a.ID)
		}
		if strings.Contains(strconv.Itoa(a.CreationDate), q) {
			add(strconv.Itoa(a.CreationDate), SuggestionCreationYear, a.ID)
		}
	}

	return suggestions
}

func containsIgnoreCase(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), sub)
}

func matchArtist(q string, a models.Artist) bool {
	if containsIgnoreCase(a.Name, q) {
		return true
	}
	for _, m := range a.Members {
		if containsIgnoreCase(m, q) {
			return true
		}
	}
	for _, loc := range a.Locations {
		if containsIgnoreCase(loc, q) {
			return true
		}
	}
	if strings.Contains(strconv.Itoa(a.CreationDate), q) {
		return true
	}
	if year, ok := firstYearFromString(a.FirstAlbum); ok && strings.Contains(strconv.Itoa(year), q) {
		return true
	}
	return false
}

var yearRegexp = regexp.MustCompile(`\b(\d{4})\b`)

func firstYearFromString(s string) (int, bool) {
	match := yearRegexp.FindStringSubmatch(s)
	if len(match) < 2 {
		return 0, false
	}
	year, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, false
	}
	return year, true
}
