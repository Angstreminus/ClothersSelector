package entity

import "github.com/google/uuid"

type Preset struct {
	Id        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	UserId    uuid.UUID `db:"user_id"`
	IsDeleted bool      `db:"is_deleted"`
	Season    string    `db:"season"`
	CreatedAt string    `db:"created_at"`
	UpdatedAt string    `db:"updated_at"`
	DeletedAt string    `db:"deleted_at"`
}
