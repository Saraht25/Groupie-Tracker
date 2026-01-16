package models

type Location struct {
    Id      int      `json:"id"`
    City    string   `json:"city"`
    Country string   `json:"country"`
}

func (l *Location) HasLocation(loc string) bool {
    return l.City == loc
}