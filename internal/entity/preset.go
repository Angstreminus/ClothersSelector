package entity

import (
	"time"
)

type Preset struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	UserId    string    `db:"user_id"`
	Season    string    `db:"season"`
	IsDeleted bool      `db:"is_deleted"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}
