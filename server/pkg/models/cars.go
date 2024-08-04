package models

import (
	"errors"

	"gorm.io/gorm"
)

var validCarTypes = map[string]bool{
	"pkw":   true,
	"lkw":   true,
	"other": true,
	"":      true,
}

type Cars struct {
	gorm.Model
	Name         string `json:"name" binding:"required"`
	Type         string `json:"type" binding:"required"`
	LicensePlate string `json:"licenseplate" binding:"required"`

	UserID *uint
	User   *User `gorm:"foreignKey:UserID"`

	OrganisationID *uint
	Organisation   *Organisation `gorm:"foreignKey:OrganisationID"`
}

func (c *Cars) BeforeCreate(tx *gorm.DB) (err error) {
	if c.UserID != nil && c.OrganisationID != nil {
		return errors.New("car cannot belong to both a user and an organization")
	}
	return nil
}

func (c *Cars) BeforeSave(tx *gorm.DB) (err error) {
	if !validCarTypes[c.Type] {
		return errors.New("invalid car value")
	}
	return nil
}
