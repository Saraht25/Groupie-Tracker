// Package models définit les structures de données du projet.
// La structure Artist est la colonne vertébrale de l'application, représentant les artistes/groupes
// avec leurs informations essentielles (nom, image, membres, dates de création, concerts, etc).
// Toutes les autres parties du code manipulent ces objets Artist.
package models

type Artist struct {
	Id              int      `json:"id"`
	Name            string   `json:"name"`
	Image           string   `json:"image"`
	Members         []string `json:"members"`
	CreationDate    int      `json:"creationDate"`
	FirstAlbum      string   `json:"firstAlbum"`
	LocationsURL    string   `json:"locations"`
	ConcertDatesURL string   `json:"concertDates"`
	RelationsURL    string   `json:"relations"`
	Locations       []string `json:"-"`
}
