package entity

import (
	"time"
)

type User struct {
	Id             string    `db:"id"`
	Login          string    `db:"login"`
	Name           string    `db:"name"`
	Surname        string    `db:"surname"`
	Role           string    `db:"role"`
	HashedPassword string    `db:"hashed_password"`
	IsDeleted      bool      `db:"is_deleted"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	DeletedAt      time.Time `db:"deleted_at"`
}
