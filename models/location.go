package models

import "strings"

type Location struct {
    Id      int      `json:"id"`
    City    string   `json:"city"`
    Country string   `json:"country"`
}

func (l *Location) HasLocation(loc string) bool {
    return l.City == loc
}

func (l *Location) MatchesQuery(query string) bool {
    q := strings.ToLower(query)          
    city := strings.ToLower(l.City)      
    country := strings.ToLower(l.Country)

    if strings.Contains(city, q) || strings.Contains(country, q){
        return true
    }
    return false
}

