package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Cars struct {
	bun.BaseModel `bun:"table:cars,alias:c"`

	ID        int64     `bun:",pk,autoincrement"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt bun.NullTime
	DeletedAt time.Time `bun:",soft_delete"`

	Name         string
	LicensePlate string

	UserID int64
}
