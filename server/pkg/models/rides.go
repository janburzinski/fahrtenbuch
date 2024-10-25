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
}

/*type Rides struct {
	bun.BaseModel `bun:"table:rides,alias:r"`

	ID        int64     `bun:",pk,autoincrement"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt bun.NullTime
	DeletedAt time.Time `bun:",soft_delete"`

	Name         string
	StartAddress string
	EndAddress   string
	Description  string
	StartTime    time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	EndTime      bun.NullTime
}
*/
