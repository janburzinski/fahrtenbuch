package models

import "gorm.io/gorm"

type Organisation struct {
	gorm.Model
	Name       string `json:"name" bindung:"required"`
	ProfilePic string `json:"profilepic" binding:"required"`

	//user
	OwnerID *uint `json:"owner_id"`
	Owner   *User `gorm:"foreignKey:UserID"`
	//TODO: Store Array of Users (Participants)
}
