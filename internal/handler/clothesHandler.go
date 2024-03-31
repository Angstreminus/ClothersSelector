package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Angstreminus/ClothersSelector/config"
	"github.com/Angstreminus/ClothersSelector/internal/apperrors"
	"github.com/Angstreminus/ClothersSelector/internal/dto"
	"github.com/Angstreminus/ClothersSelector/internal/service"
	"github.com/Angstreminus/ClothersSelector/logger"
)

type ClothesHandler struct {
	Service *service.ClothesService
	Config  *config.Config
	Logger  *logger.Logger
}

func NewClothesHandler(cfg *config.Config, log *logger.Logger, serv *service.ClothesService) *ClothesHandler {
	return &ClothesHandler{
		Service: serv,
		Config:  cfg,
		Logger:  log,
	}
}

func (ch *ClothesHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		ch.Logger.ZapLogger.Error("Error to handle request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var itemReq dto.ClthRequest
	if err = json.Unmarshal(body, &itemReq); err != nil {
		w.Header().Set("Content-Type", "application/json")
		ch.Logger.ZapLogger.Error("Error to handle request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	ch.Logger.ZapLogger.Info(r.URL.Path)
	ch.Logger.ZapLogger.Info(strings.Split(r.URL.Path, "/")[3])
	var item dto.Clothes
	item.Name = itemReq.Name
	item.Link = itemReq.Link
	item.Type = itemReq.Type
	item.PresetId = strings.Split(r.URL.Path, "/")[3]

	res, err := ch.Service.CreateItem(item)
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
