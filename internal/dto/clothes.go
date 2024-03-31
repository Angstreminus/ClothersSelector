package dto

import "time"

type ClthRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Link string `json:"link"`
}

type Clothes struct {
	Id        string    `db:"id"`
	PresetId  string    `db:"preset_id"`
	Name      string    `db:"name"`
	Type      string    `db:"type"`
	Link      string    `db:"link"`
	IsDeleted bool      `db:"id_deleted"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}
