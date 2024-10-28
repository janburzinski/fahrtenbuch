package models

import "time"

type User struct {
	BaseModel

	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`

	RefreshToken       string
	RefreshTokenExpiry time.Time

	//links multiple cars to one user
	Cars []Cars `gorm:"foreignKey:UserID"`
}
