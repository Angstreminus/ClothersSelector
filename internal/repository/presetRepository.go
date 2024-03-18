package repository

import "github.com/jmoiron/sqlx"

type PresetRepository struct {
	DB *sqlx.DB
}

func NewPresetRepository(db *sqlx.DB) *PresetRepository {
	return &PresetRepository{
		DB: db,
	}
}

func (pr *PresetRepository) CreatePreset() {
}
