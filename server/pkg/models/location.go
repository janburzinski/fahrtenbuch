package models

type Location struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	PostalCode string `json:"postalcode"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
}
