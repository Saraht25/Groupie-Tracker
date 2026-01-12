package ui

import (
	"regexp"
	"strconv"
	"strings"

	"Groupie-Tracker/models"
)

type Range struct {
	Min int
	Max int
}

type Filters struct {
	CreationYear        Range
	CreationYearEnabled bool

	FirstAlbumYear        Range
	FirstAlbumEnabled     bool

	MemberCount        Range
	MemberCountEnabled bool
	MemberCountWhitelist []int

	Locations []string
}

func ApplyFilters(f Filters, artists []models.Artist) []models.Artist {
	locationSet := buildLowerSet(f.Locations)
	memberWhitelist := buildIntSet(f.MemberCountWhitelist)

	var out []models.Artist
	for _, a := range artists {
		if !passesCreationYear(f, a) {
			continue
		}
		if !passesFirstAlbum(f, a) {
			continue
		}
		if !passesMemberCount(f, a, memberWhitelist) {
			continue
		}
		if !passesLocations(locationSet, a) {
			continue
		}
		out = append(out, a)
	}
	return out
}

func passesCreationYear(f Filters, a models.Artist) bool {
	if !f.CreationYearEnabled {
		return true
	}
	if f.CreationYear.Min != 0 && a.CreationDate < f.CreationYear.Min {
		return false
	}
	if f.CreationYear.Max != 0 && a.CreationDate > f.CreationYear.Max {
		return false
	}
	return true
}

func passesFirstAlbum(f Filters, a models.Artist) bool {
	if !f.FirstAlbumEnabled {
		return true
	}
	year, ok := firstYearFromString(a.FirstAlbum)
	if !ok {
		return false
	}
	if f.FirstAlbumYear.Min != 0 && year < f.FirstAlbumYear.Min {
		return false
	}
	if f.FirstAlbumYear.Max != 0 && year > f.FirstAlbumYear.Max {
		return false
	}
	return true
}

func passesMemberCount(f Filters, a models.Artist, whitelist map[int]struct{}) bool {
	count := len(a.Members)

	if len(whitelist) > 0 {
		_, ok := whitelist[count]
		return ok
	}

	if !f.MemberCountEnabled {
		return true
	}
	if f.MemberCount.Min != 0 && count < f.MemberCount.Min {
		return false
	}
	if f.MemberCount.Max != 0 && count > f.MemberCount.Max {
		return false
	}
	return true
}

func passesLocations(locationSet map[string]struct{}, a models.Artist) bool {
	if len(locationSet) == 0 {
		return true
	}
	for _, loc := range a.Locations {
		if _, ok := locationSet[strings.ToLower(loc)]; ok {
			return true
		}
	}
	return false
}

func buildLowerSet(values []string) map[string]struct{} {
	if len(values) == 0 {
		return nil
	}
	set := make(map[string]struct{}, len(values))
	for _, v := range values {
		set[strings.ToLower(strings.TrimSpace(v))] = struct{}{}
	}
	return set
}

func buildIntSet(values []int) map[int]struct{} {
	if len(values) == 0 {
		return nil
	}
	set := make(map[int]struct{}, len(values))
	for _, v := range values {
		set[v] = struct{}{}
	}
	return set
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
	return year, tru
}
