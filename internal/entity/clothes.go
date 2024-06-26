package entity

import "time"

type Clothes struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	Type      string    `db:"type"`
	Link      string    `db:"link"`
	IsDeleted bool      `db:"id_deleted"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}
