package models

type Relation struct {
    ArtistId   int `json:"artistId"`
    DateId     int `json:"dateId"`
    LocationId int `json:"locationId"`
}

func (r *Relation) MatchesArtist(id int) bool {
	return id == r.ArtistId
}

func (r *Relation) MatchesDate(id int) bool {
	return id == r.DateId
}

func (r *Relation) MatchesLocation(id int) bool {
	return id == r.LocationId
}