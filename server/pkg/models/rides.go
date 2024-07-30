package models

import "gorm.io/gorm"

type Rides struct {
	gorm.Model
	Title       string   `json:"title" binding:"required"`
	Description *string  `json:"description"`
	RideFrom    *string  `json:"ride_from"`
	RideTo      *string  `json:"ride_to"`
	Stops       []string `json:"stops" gorm:"type:text[]"`
	Distance    *string  `json:"distance"`
	BeginTime   string   `json:"begin_time"`
	EndTime     string   `json:"end_time"`
	Category    string   `json:"category"` // something like work, shopping, freetime etc...

	//user
	UserID *uint `json:"user_id"`
	User   *User `gorm:"foreignKey:UserID"`

	//organisation
	OrganisationID *uint         `json:"organisation_id"`
	Organisation   *Organisation `gorm:"foreignkey:OrganisationID"`
}
