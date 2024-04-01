package repository

import (
	"time"

	"github.com/Angstreminus/ClothersSelector/internal/apperrors"
	"github.com/Angstreminus/ClothersSelector/internal/dto"
	"github.com/Angstreminus/ClothersSelector/internal/entity"
	"github.com/Angstreminus/ClothersSelector/logger"
	"github.com/jmoiron/sqlx"
)

type ClothesRepository struct {
	Db     *sqlx.DB
	Logger *logger.Logger
}

func (cr *ClothesRepository) Create(item dto.Clothes) (*entity.Clothes, apperrors.AppError) {
	item.CreatedAt = time.Now().Local().UTC()
	query := "INSERT INTO clothes (id, name, season, user_id, hashed_password, is_deleted, created_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, login, name, surname, role, hashed_password, is_deleted, created_at;"
	var cloth entity.Clothes
	tx, err := cr.Db.Beginx()
	if err != nil {
		cr.Logger.ZapLogger.Error("Db Tx begin error")
		return nil, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	row := tx.QueryRow(query, &item.Id, &item.Name, &item.Type, &item.Link, &item.IsDeleted, &item.CreatedAt)
	if err := row.Scan(&cloth); err != nil {
		tx.Rollback()
		cr.Logger.ZapLogger.Error("Query error")
		return nil, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}

	queryCP := "INSERT INTO clothers_presets (preset_id, cloth_id) VALUES($1, $2);"
	res, err := tx.Exec(queryCP, &item.Id, &item.Id)
	if err != nil {
		tx.Rollback()
		cr.Logger.ZapLogger.Error("Query error")
		return nil, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	rowsAff, err := res.RowsAffected()
	if rowsAff == 0 || err != nil {
		tx.Rollback()
		cr.Logger.ZapLogger.Error("Clothes presets not created")
		return nil, &apperrors.DBoperationErr{
			Message: "Clothes presets not created",
		}
	}
	tx.Commit()
	cr.Logger.ZapLogger.Info("Cloth item created")
	return &cloth, nil
}

func (cr *ClothesRepository) GetClothes(clothId, presetId, season, userId string) ([]entity.Clothes, apperrors.AppError) {
	query := `SELECT c.id, c.name, c.type, c.link FROM clothes c JOIN clothes_presets cp ON $1 = cp.cloth_id
	JOIN presets p ON cp.preset_id = $2
	WHERE p.season = $3 AND p.user_id = $4;`
	var clothes []entity.Clothes
	rows, err := cr.Db.Queryx(query, &clothId, &presetId, &season, &userId)
	if err != nil {
		cr.Logger.ZapLogger.Error("Getting clothes quering error")
		return nil, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	for rows.Next() {
		err := rows.StructScan(&clothes)
		if err != nil {
			cr.Logger.ZapLogger.Error("Getting clothes quering error")
			return nil, &apperrors.DBoperationErr{
				Message: err.Error(),
			}
		}
	}
	return clothes, nil
}
