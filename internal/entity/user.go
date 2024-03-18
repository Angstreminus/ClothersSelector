package entity

import "github.com/google/uuid"

type User struct {
	Id             uuid.UUID `db:"id"`
	Login          string    `db:"login"`
	Name           string    `db:"name"`
	Surname        string    `db:"surname"`
	Role           string    `db:"role"`
	HashedPassword string    `db:"hashed_password"`
	IsDeleted      bool      `db:"is_deleted"`
	CreatedAt      string    `db:"created_at"`
	UpdatedAt      string    `db:"updated_at"`
	DeletedAt      string    `db:"deleted_at"`
}
