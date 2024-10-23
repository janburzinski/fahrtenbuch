package models

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64     `bun:",pk,autoincrement"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt bun.NullTime
	DeletedAt time.Time `bun:",soft_delete"`

	FirstName string
	LastName  string
	Email     string `bun:",unique"`
	Password  string

	Cars []*Cars `bun:"rel:has-many,join:id=user_id"`
}
