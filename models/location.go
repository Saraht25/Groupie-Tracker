package models

type Location struct {
    Id      int      `json:"id"`
    City    string   `json:"city"`
    Country string   `json:"country"`
}