package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Angstreminus/ClothersSelector/internal/apperrors"
	"github.com/Angstreminus/ClothersSelector/internal/chache"
	"github.com/Angstreminus/ClothersSelector/internal/dto"
	"github.com/Angstreminus/ClothersSelector/internal/entity"
	"github.com/Angstreminus/ClothersSelector/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB     *sqlx.DB
	Chache *chache.Chache
	Logger *logger.Logger
}

func NewUserRepository(chache *chache.Chache, db *sqlx.DB, log *logger.Logger) *UserRepository {
	return &UserRepository{
		Chache: chache,
		DB:     db,
		Logger: log,
	}
}

func (ur *UserRepository) RegisterUser(user *entity.User, ctx *context.Context) (*entity.User, apperrors.AppError) {
	user.CreatedAt = time.Now().Format("2006-01-02 15:04:05.999")
	user.IsDeleted = false
	query := `INSERT INTO users VALUES(id, login, name, surname, hashed_password, is_deleted) $1, $2, $3, $4, $5, $6 RETURNING *;`
	stmt, err := ur.DB.PrepareContext(*ctx, query)
	if err != nil {
		ur.Logger.ZapLogger.Error("Preparing statement error")
		return nil, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	var usr entity.User
	err = stmt.QueryRowContext(*ctx, &user).Scan(&usr)
	if err != nil {
		ur.Logger.ZapLogger.Error("Query error")
		return nil, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	ur.Logger.ZapLogger.Info("User created")
	return &usr, nil
}

func (ur *UserRepository) UpdateUser(ctx *context.Context, toUpdate *entity.User) (*entity.User, apperrors.AppError) {
	toUpdate.UpdatedAt = time.Now().Format("2006-01-02 15:04:05.999")
	toUpdate.IsDeleted = false
	query := `
	UPDATE users
	SET name=$1, login=#2, surname=$3, hashed_password=$4, is_deleted=$5, updated_at=$6)
	WHERE id::text=$6 RETURNINIG *;`
	var res entity.User
	stmt, err := ur.DB.PrepareContext(*ctx, query)
	ur.Logger.ZapLogger.Error("")
	if err != nil {
		ur.Logger.ZapLogger.Error("Preparing statement error")
		return nil, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	if err = stmt.QueryRowContext(*ctx, &toUpdate).Scan(res); err != nil {
		if err == sql.ErrNoRows {
			ur.Logger.ZapLogger.Error("Update: No rows error")
			return nil, &apperrors.DBoperationErr{
				Message: err.Error(),
			}
		}
		return nil, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	ur.Logger.ZapLogger.Info("User created")
	return &res, nil
}

func (ur *UserRepository) DeleteUser(ctx *context.Context, userId uuid.UUID) apperrors.AppError {
	deletedAt := time.Now().Format("2006-01-02 15:04:05.999")
	isDeleted := true
	query := `
	UPDATE users
	SET is_deleted=$1, deleted_at=$2)
	WHERE id::text=$3;`
	stmt, err := ur.DB.PrepareContext(*ctx, query)
	if err != nil {
		ur.Logger.ZapLogger.Error("Preparing statement error")
		return &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	row, err := stmt.Exec(isDeleted, deletedAt, userId)
	if err != nil {
		return &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	rowsAff, err := row.RowsAffected()
	if err != nil {
		return &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	if rowsAff == 0 {
		return &apperrors.DBoperationErr{
			Message: "No rows affected",
		}
	}
	ur.Logger.ZapLogger.Info("User successfully deleted")
	return nil
}

func (ur *UserRepository) GetUser(ctx *context.Context, userId uuid.UUID) (*entity.User, apperrors.AppError) {
	query := `SELECT * FROM users WHERE id::text=$1;`
	var res *entity.User
	stmt, err := ur.DB.PrepareContext(*ctx, query)
	if err != nil {
		ur.Logger.ZapLogger.Error("Preparing statement error")
		return nil, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	if err = stmt.QueryRowContext(*ctx, userId).Scan(res); err != nil {
		ur.Logger.ZapLogger.Error("Error while scanning entity")
		return nil, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	ur.Logger.ZapLogger.Info("User extracted successfully")
	return res, nil
}

func (ur *UserRepository) UserIsExist(ctx *context.Context, userSignature *dto.UserSignature) (bool, error) {
	query := `SELECT EXIST(SELECT 1 FROM users WHERE login = $1 AND hashed_password = $2);`
	stmt, err := ur.DB.PrepareContext(*ctx, query)
	if err != nil {
		ur.Logger.ZapLogger.Error("Preparing statement error")
		return true, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	var exist bool
	err = stmt.QueryRowContext(*ctx, &userSignature.Login, &userSignature.Password).Scan(&exist)
	if err != nil {
		ur.Logger.ZapLogger.Error("Query error")
		return true, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	ur.Logger.ZapLogger.Info("User does not exist")
	return exist, nil
}
