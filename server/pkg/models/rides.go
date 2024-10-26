package models

import (
	"time"
)

type Rides struct {
	BaseModel

	Name         string `gorm:"not null"`
	StartAddress string `gorm:"not null"`
	EndAddress   string
	Description  string
	StarTime     time.Time `gorm:"not null"`
	EndTime      time.Time

	CarID uint
}
