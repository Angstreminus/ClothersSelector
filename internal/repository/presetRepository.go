package repository

import (
	"database/sql"
	"errors"
	"fmt"
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
	query := "INSERT INTO presets (id, name, season, user_id, is_deleted, created_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, name, season, user_id, is_deleted, created_at;"
	var prst entity.Preset
	fmt.Println("Preset userId")
	fmt.Println(preset.UserId)
	fmt.Println("Preset userId end")
	row := pr.DB.QueryRowx(query, &preset.Id, &preset.Name, &preset.Season, &preset.UserId, &preset.IsDeleted, &preset.CreatedAt)
	if err := row.StructScan(&prst); err != nil {
		pr.Logger.ZapLogger.Error("Query error")
		return nil, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	pr.Logger.ZapLogger.Info("Preset created")
	return &prst, nil
}

func (pr *PresetRepository) GetPresetById(id string) (*entity.Preset, apperrors.AppError) {
	query := `SELECT id, name, season, user_id, is_deleted, created_at FROM presets WHERE (id=$1);`
	row := pr.DB.QueryRowx(query, &id)
	var preset entity.Preset
	if err := row.StructScan(&preset); err != nil {
		pr.Logger.ZapLogger.Error("Preparing statement error")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apperrors.AuthError{
				Message: "Such user does not exists",
			}
		}
		return nil, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	pr.Logger.ZapLogger.Info("User found successfully")
	return &preset, nil
}
