package models

import (
	"strings"
	"strconv"
)

type Date struct {
    Id    int      `json:"id"`
    Dates []string `json:"dates"`
}

func (d *Date) HasDate(date string) bool{
	for _ , dateStr := range d.Dates{
		if dateStr == date{
			return true
		}
	}
	return false
}	


func (d *Date) MatchesQuery(query string) bool {
	query = strings.ToLower(query)

	for _, dateStr := range d.Dates {
    	if strings.Contains(strings.ToLower(dateStr), query) {
            return true
        }
    }
    return false
}

func (d *Date) ExtractYear(date string) int {
    if len(date) < 4 {
        return 0
    }

    yearStr := date[:4]
    year, err := strconv.Atoi(yearStr)
    if err != nil {
        return 0
    }

    return year
}

func (d *Date) HasDateInRange(minYear, maxYear int) bool {
	for _ , datestr := range d.Dates{
		if d.ExtractYear(datestr) >= minYear && d.ExtractYear(datestr) <= maxYear{
			return true
		}
	}
	return false
}