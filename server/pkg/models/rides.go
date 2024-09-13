package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Rides struct {
	gorm.Model
	Title       string   `json:"title" binding:"required"`
	Description *string  `json:"description"`
	RideFrom    Location `json:"ride_from" gorm:"embedded;embeddedPrefix:from_"`
	RideTo      Location `json:"ride_to" gorm:"embedded;embeddedPrefix:from_"`
	Stops       []Stop   `json:"stops" gorm:"serializer:json"`
	//	Route       Route    `json:"route" gorm:"embedded"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Category  string `json:"category"` // something like work, shopping, freetime etc...

	// car
	// link ride to car
	CarID uint `json:"car_id" binding:"required"`
	Car   Cars `gorm:"foreignKey:CarID"`

	//user
	// save who created the ride
	UserID *uint `json:"user_id" binding:"required"`
	User   *User `gorm:"foreignKey:UserID"`

	// //organisation
	// OrganisationID *uint         `json:"organisation_id"`
	// Organisation   *Organisation `gorm:"foreignkey:OrganisationID"`
}

func (r *Rides) Validate() error {
	if r.Title == "" {
		return errors.New("title is required")
	}

	if r.UserID == nil || *r.UserID == 0 {
		return errors.New("userid is not set")
	}

	if r.CarID == 0 {
		return errors.New("car id is not set")
	}

	if err := r.ValidateTimes(); err != nil {
		return err
	}

	return nil
}

func (r *Rides) ValidateTimes() error {
	layout := "2006-01-02 15:04:05"
	start, err := time.Parse(layout, r.StartTime)
	if err != nil {
		return errors.New("invalid start time format")
	}

	end, err := time.Parse(layout, r.EndTime)
	if err != nil {
		return errors.New("invalid end time format")
	}

	if end.Before(start) {
		return errors.New("end time is before start time")
	}

	return nil
}

func (r *Rides) BeforeCreate(tx *gorm.DB) error {
	return r.Validate()
}

func (r *Rides) BeforeUpdate(tx *gorm.DB) error {
	return r.Validate()
}
