package services

import "Groupie-Tracker/models"

func FilterArtists(
	artists []models.Artist,
	dates []models.Date,
	locations []models.Location,
	relations []models.Relation,
	minYear int,
	maxYear int,
	city string,
	country string,
) []models.Artist {

	var result []models.Artist

	for _, artist := range artists {

		match := true

		// 🔹 Filtre par date
		if minYear != 0 || maxYear != 0 {
			hasDate := false

			for _, rel := range relations {
				if !rel.MatchesArtist(artist.Id) {
					continue
				}

				for _, d := range dates {
					if rel.MatchesDate(d.Id) && d.HasDateInRange(minYear, maxYear) {
						hasDate = true
						break
					}
				}
			}

			if !hasDate {
				match = false
			}
		}

		// 🔹 Filtre par location
		if match && (city != "" || country != "") {
			hasLocation := false

			for _, rel := range relations {
				if !rel.MatchesArtist(artist.Id) {
					continue
				}

				for _, l := range locations {
					if rel.MatchesLocation(l.Id) {

						if city != "" && l.City == city {
							hasLocation = true
						}

						if country != "" && l.Country == country {
							hasLocation = true
						}
					}
				}
			}

			if !hasLocation {
				match = false
			}
		}

		if match {
			result = append(result, artist)
		}
	}

	return result
}
