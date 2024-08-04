package models

import "gorm.io/gorm"

type Organisation struct {
	gorm.Model
	Name       string `json:"name" bindung:"required"`
	ProfilePic string `json:"profilepic" binding:"required"`

	//user
	OwnerID      uint
	Owner        User   `gorm:"foreignKey:OwnerID"`
	Participants []User `gorm:"foreignKey:OrganisationID"`

	//cars
	Cars []Cars `gorm:"foreignKey:OrganisationID"`
}
