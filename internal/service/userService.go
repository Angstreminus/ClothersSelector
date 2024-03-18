package service

import (
	"github.com/Angstreminus/ClothersSelector/internal/apperrors"
	"github.com/Angstreminus/ClothersSelector/internal/dto"
	"github.com/Angstreminus/ClothersSelector/internal/entity"
	"github.com/Angstreminus/ClothersSelector/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	Ur *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		Ur: repo,
	}
}

func (us *UserService) RegiterUser(toRegistrate *dto.ReisterRequest) (*entity.User, apperrors.AppError) {
	var (
		user entity.User
		err  error
	)
	user.Id = uuid.New()
	user.Name = toRegistrate.Name
	user.Login = toRegistrate.Login
	user.Surname = toRegistrate.Surname
	user.HashedPassword, err = HashPassword(toRegistrate.Password)
	if err != nil {
		us.Ur.Logger.ZapLogger.Error("Error to hash password")
		return nil, &apperrors.HashError{
			Message: err.Error(),
		}
	}
	user.Role = "User"
	user.IsDeleted = false
	return &user, nil
}
