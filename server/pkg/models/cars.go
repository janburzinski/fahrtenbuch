package models

type Cars struct {
	BaseModel

	Name         string `gorm:"not null"`
	LicensePlate string

	UserID int
}
