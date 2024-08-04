package models

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

var validRanks = map[string]bool{
	"owner":  true,
	"editor": true,
	"user":   true,
	"":       true,
}

type User struct {
	gorm.Model
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required" gorm:"unique;not null"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required" gorm:"not null"`

	//organisation
	OrganisationID      *uint
	Organisation        *Organisation `gorm:"foreignKey:OrganisationID;references:ID"`
	Rank                string
	OwnedOrganisationID *uint
	OwnedOrganisation   *Organisation `gorm:"foreignKey:OwnedOrganisationID;references:ID"`

	//cars
	Cars []Cars `gorm:"foreignKey:UserID"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	//make use that the rank value is something valid
	if !validRanks[u.Rank] {
		log.Printf("Error: Invalid rank value %s", u.Rank)
		return errors.New("invalid rank value")
	}
	return nil
}
