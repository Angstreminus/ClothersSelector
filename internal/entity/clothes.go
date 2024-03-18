package entity

import "github.com/google/uuid"

type Clothes struct {
	Id        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Type      string    `db:"type"`
	Link      string    `db:"link"`
	IsDeleted bool      `db:"id_deleted"`
	CreatedAt string    `db:"created_at"`
	UpdatedAt string    `db:"updated_at"`
	DeletedAt string    `db:"deleted_at"`
}
