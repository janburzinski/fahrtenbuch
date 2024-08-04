package models

type Route struct {
	EncodedPolyline string `json:"encoded_polyline"`
	Distance        int    `json:"distance_meter"`
	Duration        int    `json:"duration_seconds"`
}
