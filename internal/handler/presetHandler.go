package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Angstreminus/ClothersSelector/config"
	"github.com/Angstreminus/ClothersSelector/internal/apperrors"
	"github.com/Angstreminus/ClothersSelector/internal/dto"
	"github.com/Angstreminus/ClothersSelector/internal/entity"
	"github.com/Angstreminus/ClothersSelector/internal/service"
	"github.com/Angstreminus/ClothersSelector/logger"
)

type PresetHandler struct {
	Serv   *service.PresetService
	Logger *logger.Logger
	Config *config.Config
}

func NewPresetHandler(cfg *config.Config, serv *service.PresetService, log *logger.Logger) *PresetHandler {
	return &PresetHandler{
		Serv:   serv,
		Config: cfg,
		Logger: log,
	}
}

func (ph *PresetHandler) CreatePreset(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		ph.Logger.ZapLogger.Error("Error to handle request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var prReq dto.PresetCreateRequest
	if err = json.Unmarshal(body, &prReq); err != nil {
		w.Header().Set("Content-Type", "application/json")
		ph.Logger.ZapLogger.Error("Error to handle request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	ph.Logger.ZapLogger.Info(r.URL.Path)
	ph.Logger.ZapLogger.Info(strings.Split(r.URL.Path, "/")[1])
	var preset entity.Preset
	preset.Name = prReq.Name
	preset.Season = prReq.Season
	preset.UserId = strings.Split(r.URL.Path, "/")[1]

	res, err := ph.Serv.CreatePreset(preset)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		respErr := apperrors.MatchError(err)
		w.WriteHeader(respErr.Status)
		_ = json.NewEncoder(w).Encode(respErr)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
