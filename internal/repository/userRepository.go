package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (ur *UserRepository) RegisterUser(user entity.User) (*entity.User, apperrors.AppError) {
	user.CreatedAt = time.Now().Local().UTC()
	query := "INSERT INTO users (id, login, name, surname, role, hashed_password, is_deleted, created_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, login, name, surname, role, hashed_password, is_deleted, created_at;"
	var usr entity.User
	row := ur.DB.QueryRowx(query, &user.Id, &user.Login, &user.Name, &user.Surname, &user.Role, &user.HashedPassword, &user.IsDeleted, &user.CreatedAt)
	if err := row.StructScan(&usr); err != nil {
		ur.Logger.ZapLogger.Error("Query error")
		return nil, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	ur.Logger.ZapLogger.Info("User created")
	return &usr, nil
}

func (ur *UserRepository) UserExists(login string) (bool, apperrors.AppError) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE login=$1);`
	var exist bool = false
	err := ur.DB.QueryRowx(query, login).Scan(&exist)
	if err != nil {
		if err == sql.ErrNoRows {
			ur.Logger.ZapLogger.Info("User does not exist")
			return false, nil
		}
		ur.Logger.ZapLogger.Error("Existance Query error")
		return true, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	return exist, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, toUpdate *entity.User) (*entity.User, apperrors.AppError) {
	toUpdate.UpdatedAt = time.Now().Local().UTC()
	toUpdate.IsDeleted = false
	query := `
	UPDATE users
	SET name=$1, login=#2, surname=$3, hashed_password=$4, is_deleted=$5, updated_at=$6)
	WHERE id::text=$6 RETURNINIG *;`
	var res entity.User
	stmt, err := ur.DB.PrepareContext(ctx, query)
	ur.Logger.ZapLogger.Error("")
	if err != nil {
		ur.Logger.ZapLogger.Error("Preparing statement error")
		return nil, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	if err = stmt.QueryRowContext(ctx, &toUpdate).Scan(res); err != nil {
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

func (ur *UserRepository) GetUserByLogin(login dto.LoginRequest) (*entity.User, apperrors.AppError) {
	query := `SELECT id, login, name, surname, role, hashed_password, is_deleted FROM users WHERE (login=$1);`
	fmt.Println("Repository started log in operation")
	row := ur.DB.QueryRowx(query, &login.Login)
	var user entity.User
	if err := row.StructScan(&user); err != nil {
		ur.Logger.ZapLogger.Error("Preparing statement error")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apperrors.AuthError{
				Message: "Such user does not exists",
			}
		}
		return nil, &apperrors.DBoperationErr{
			Message: err.Error(),
		}
	}
	ur.Logger.ZapLogger.Info("User found successfully")
	return &user, nil
}
