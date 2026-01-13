package models

import (
	"strings"
	"fmt"
	"strconv"
)
		

type Artist struct {
    Id          int      `json:"id"`
    Name        string   `json:"name"`
    Image       string   `json:"image"`
    Members     []string `json:"members"`
    CreationDate int     `json:"creationDate"`
    FirstAlbum  string   `json:"firstAlbum"`
}

func (a *Artist) GetMembersCount() int {
	return len(a.Members)
}

func (a *Artist) HasMember(name string) bool {
	for _ , member := range a.Members {
		if member == name {
			return true
		}
	}
	return false
}

func (a *Artist) MatchesSearch(query string) bool{
	query = strings.ToLower(query)
	// Vérifier le nom de l'artiste
	if strings.Contains(strings.ToLower(a.Name), query) {
		return true
	}

	// Vérifier chaque membre
	for _, member := range a.Members {
		if strings.Contains(strings.ToLower(member), query) {
			return true
		}
	}

	// Vérifier le premier album
	if strings.Contains(strings.ToLower(a.FirstAlbum), query) {
		return true
	}

	// Vérifier la date de création (convertir en string)
	if strings.Contains(strings.ToLower(fmt.Sprint(a.CreationDate)), query) {
		return true
	}

	// Si aucun champ ne correspond
	return false
}


func (a *Artist) IsCreatedBetween(start, end int) bool{
	return a.CreationDate >= start && a.CreationDate <= end
}



func (a *Artist) isFirstAlbumBetween(start, end int) bool{
	album, err := strconv.Atoi(a.FirstAlbum)
	if err != nil{
		return false
	}
	return album >= start && album <= end
}