package models

import (
	"errors"

	"gorm.io/gorm"
)

const (
	RankOwner  = "owner"
	RankEditor = "editor"
	RankUser   = "user"
)

type User struct {
	gorm.Model
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required" gorm:"unique,not null"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`

	//organisation
	OrganisationID *uint         `json:"organisation_id"`
	Organisation   *Organisation `gorm:"foreignkey:OrganisationID"`
	Rank           string        `json:"rank"` // used for permission inside the organisation
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	//make use that the rank value is something valid
	if u.Rank != RankOwner && u.Rank != RankEditor && u.Rank != RankUser && u.Rank != "" {
		return errors.New("invalid rank value")
	}
	return nil
}
