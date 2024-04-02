package service

import (
	"github.com/Angstreminus/ClothersSelector/internal/apperrors"
	"github.com/Angstreminus/ClothersSelector/internal/dto"
	"github.com/Angstreminus/ClothersSelector/internal/entity"
	"github.com/Angstreminus/ClothersSelector/internal/repository"
	"github.com/Angstreminus/ClothersSelector/logger"
	"github.com/google/uuid"
)

type ClothesService struct {
	Logger *logger.Logger
	Repo   *repository.ClothesRepository
}

func NewClothesService(repo *repository.ClothesRepository, log *logger.Logger) *ClothesService {
	return &ClothesService{
		Repo:   repo,
		Logger: log,
	}
}

func (cs *ClothesService) CreateItem(toCreate dto.Clothes) (entity.Clothes, apperrors.AppError) {
	toCreate.Id = uuid.New().String()
	toCreate.IsDeleted = false
	return cs.Repo.Create(toCreate)
}

func (cs *ClothesService) GetItems(clothId, presetId, season, userId string) ([]entity.Clothes, apperrors.AppError) {
	return cs.Repo.GetClothes(clothId, presetId, season, userId)
}
