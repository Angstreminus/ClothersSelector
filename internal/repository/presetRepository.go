package repository

import (
	"time"

	"github.com/Angstreminus/ClothersSelector/internal/apperrors"
	"github.com/Angstreminus/ClothersSelector/internal/entity"
	"github.com/Angstreminus/ClothersSelector/logger"
	"github.com/jmoiron/sqlx"
)

type PresetRepository struct {
	DB     *sqlx.DB
	Logger *logger.Logger
}

func NewPresetRepository(db *sqlx.DB, log *logger.Logger) *PresetRepository {
	return &PresetRepository{
		DB:     db,
		Logger: log,
	}
}

func (pr *PresetRepository) CreatePreset(preset *entity.Preset) (*entity.Preset, apperrors.AppError) {
	preset.CreatedAt = time.Now().Local().UTC()
	query := "INSERT INTO presets (id, name, season, user_id, hashed_password, is_deleted, created_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, login, name, surname, role, hashed_password, is_deleted, created_at;"
	var prst entity.Preset
	row := pr.DB.QueryRowx(query, &preset.UserId, &preset.Season, &preset.IsDeleted, &preset.CreatedAt)
	if err := row.StructScan(&prst); err != nil {
		pr.Logger.ZapLogger.Error("Query error")
		return nil, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	pr.Logger.ZapLogger.Info("Preset created")
	return &prst, nil
}
