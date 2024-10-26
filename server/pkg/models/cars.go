package models

type Cars struct {
	BaseModel

	Name         string `gorm:"not null"`
	LicensePlate string

	//link (possibly) multiple cars to one user
	UserID int

	//link multiple rides to one car
	Rides []Rides `gorm:"foreignKey:CarID"`
}
