package service

import (
	"github.com/Angstreminus/ClothersSelector/internal/apperrors"
	"github.com/Angstreminus/ClothersSelector/internal/entity"
	"github.com/Angstreminus/ClothersSelector/internal/repository"
	"github.com/Angstreminus/ClothersSelector/logger"
	"github.com/google/uuid"
)

type PresetService struct {
	Logger *logger.Logger
	Repo   *repository.PresetRepository
}

func NewPresetService(repo *repository.PresetRepository, log *logger.Logger) *PresetService {
	return &PresetService{
		Logger: log,
		Repo:   repo,
	}
}

func (ps *PresetService) CreatePreset(toCreate entity.Preset) (*entity.Preset, apperrors.AppError) {
	toCreate.Id = uuid.New().String()
	toCreate.IsDeleted = false
	return ps.Repo.CreatePreset(&toCreate)
}
