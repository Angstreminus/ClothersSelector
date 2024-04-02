package service

import (
	"github.com/Angstreminus/ClothersSelector/internal/apperrors"
	"github.com/Angstreminus/ClothersSelector/internal/dto"
	"github.com/Angstreminus/ClothersSelector/internal/entity"
	"github.com/Angstreminus/ClothersSelector/internal/repository"
	"github.com/Angstreminus/ClothersSelector/logger"
	"github.com/google/uuid"
)

type UserService struct {
	Ur     *repository.UserRepository
	Logger *logger.Logger
}

func NewUserService(repo *repository.UserRepository, log *logger.Logger) *UserService {
	return &UserService{
		Ur:     repo,
		Logger: log,
	}
}

func (us *UserService) RegisterUser(toRegistrate dto.RegisterRequest) (entity.User, apperrors.AppError) {
	hashedPassword := HashPassword(toRegistrate.Password)
	exists, err := us.UserExists(toRegistrate.Login)
	if err != nil {
		us.Logger.ZapLogger.Error("User already exists")
		return entity.User{}, err
	}
	if !exists {
		us.Logger.ZapLogger.Info("user does not exists")
		var usr entity.User = entity.User{
			Name:           toRegistrate.Name,
			Surname:        toRegistrate.Surname,
			Login:          toRegistrate.Login,
			Id:             uuid.New().String(),
			Role:           "User",
			IsDeleted:      false,
			HashedPassword: string(hashedPassword),
		}
		us.Logger.ZapLogger.Info("User filled in serveice")
		return us.Ur.RegisterUser(usr)
	}
	return entity.User{}, &apperrors.UserExistError{
		Message: "User already exists",
	}
}

func (us *UserService) LoginUser(loginReq dto.LoginRequest) (*entity.User, apperrors.AppError) {
	user, err := us.Ur.GetUserByLogin(loginReq)
	if err != nil {
		return nil, err
	}

	if (loginReq.Login == user.Login) && CompareToHash(user.HashedPassword, loginReq.Password) && !user.IsDeleted {
		return nil, &apperrors.AuthError{
			Message: "Login and password does not match",
		}
	}
	return user, nil
}

func (us *UserService) UserExists(signature string) (bool, apperrors.AppError) {
	return us.Ur.UserExists(signature)
}
