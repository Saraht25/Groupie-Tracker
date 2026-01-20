package services

import "Groupie-Tracker/models"

func SearchArtists(
	artists []models.Artist,
	dates []models.Date,
	locations []models.Location,
	relations []models.Relation,
	query string,
) []models.Artist {

	var result []models.Artist
	for _, artist := range artists {
		if artist.MatchesSearch(query) {
			result = append(result, artist)
			continue
		}

		found := false

		for _, rel := range relations {
			if !rel.MatchesArtist(artist.Id) {
				continue
			}
			// Dates
			for _, d := range dates {
				if rel.MatchesDate(d.Id) && d.MatchesQuery(query) {
					found = true
					break
				}
			}
			// Locations
			for _, l := range locations {
				if rel.MatchesLocation(l.Id) && l.MatchesQuery(query) {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if found {
			result = append(result, artist)
		}
	}
	return result
}
