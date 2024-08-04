package models

import "gorm.io/gorm"

type Rides struct {
	gorm.Model
	Title       string   `json:"title" binding:"required"`
	Description *string  `json:"description"`
	RideFrom    Location `json:"ride_from" gorm:"embedded;embeddedPrefix:from_"`
	RideTo      Location `json:"ride_to" gorm:"embedded;embeddedPrefix:from_"`
	Stops       []Stop   `json:"stops" gorm:"serializer:json"` //TODO: Make it possible to store google maps routes
	Route       Route    `json:"route" gorm:"embedded"`
	BeginTime   string   `json:"begin_time"`
	EndTime     string   `json:"end_time"`
	Category    string   `json:"category"` // something like work, shopping, freetime etc...

	// car
	CarID uint `json:"car_id" binding:"required"`
	Car   Cars `gorm:"foreignKey:CarID"` // TODO: Add TESTs

	//user
	// UserID *uint `json:"user_id"`
	// User   *User `gorm:"foreignKey:UserID"`

	// //organisation
	// OrganisationID *uint         `json:"organisation_id"`
	// Organisation   *Organisation `gorm:"foreignkey:OrganisationID"`
}
