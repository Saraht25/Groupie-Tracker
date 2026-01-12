package models

type Artist struct {
	ID            int                 `json:"id"`
	Image         string              `json:"image"`
	Name          string              `json:"name"`
	Members       []string            `json:"members"`
	CreationDate  int                 `json:"creationDate"`
	FirstAlbum    string              `json:"firstAlbum"`
	Locations     []string            `json:"locations,omitempty"`
	Dates         []string            `json:"dates,omitempty"`
	Relations     map[string][]string `json:"relations,omitempty"`
	SpotifyURL    string              `json:"spotifyURL,omitempty"`
	AdditionalURL string              `json:"additionalURL,omitempty"`
}
