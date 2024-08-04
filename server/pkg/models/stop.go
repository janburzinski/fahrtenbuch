package models

type Stop struct {
	Location
	Order int `json:"order"`
}
