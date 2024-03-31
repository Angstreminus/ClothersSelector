package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Angstreminus/ClothersSelector/config"
	"github.com/Angstreminus/ClothersSelector/internal/apperrors"
	"github.com/Angstreminus/ClothersSelector/internal/auth/token"
	"github.com/Angstreminus/ClothersSelector/internal/dto"
	"github.com/Angstreminus/ClothersSelector/internal/service"
	"github.com/Angstreminus/ClothersSelector/logger"
)

type Cookie struct {
	Refresh string `json:"login" validate:"required"`
}

type UserHandler struct {
	UserService *service.UserService
	Logger      *logger.Logger
	Config      *config.Config
}

func NewUserHandler(usrServ *service.UserService, log *logger.Logger, cfg *config.Config) *UserHandler {
	return &UserHandler{
		UserService: usrServ,
		Logger:      log,
		Config:      cfg,
	}
}

func (uh *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		uh.Logger.ZapLogger.Error("Error to handle request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var toRegistrate dto.RegisterRequest
	if err = json.Unmarshal(body, &toRegistrate); err != nil {
		w.Header().Set("Content-Type", "application/json")
		uh.Logger.ZapLogger.Error("Error to handle request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	user, err := uh.UserService.RegisterUser(toRegistrate)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		respErr := apperrors.MatchError(err)
		w.WriteHeader(respErr.Status)
		_ = json.NewEncoder(w).Encode(respErr)
		return
	}
	tokenPair, err := token.CreateTokenPair(user, uh.Config)
	if err != nil {
		respErr := apperrors.MatchError(err)
		w.WriteHeader(respErr.Status)
		_ = json.NewEncoder(w).Encode(respErr)
		return
	}
	ctx := context.Background()
	exp, err := strconv.Atoi(uh.Config.RefExp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = uh.UserService.Ur.Chache.RedisChahe.Set(ctx, tokenPair.Refresh, user.Id, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		uh.Logger.ZapLogger.Error("Error to save token")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Enable to save token"))
		return
	}
	cookie := http.Cookie{
		Name:     "Refresh",
		Value:    tokenPair.Refresh,
		Expires:  time.Now().Add(time.Duration(exp) * time.Minute),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	w.Header().Set("Authorization", "Bearer "+tokenPair.Access)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (uh *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		uh.Logger.ZapLogger.Error("Error to handle request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var loginReq dto.LoginRequest
	if err = json.Unmarshal(body, &loginReq); err != nil {
		w.Header().Set("Content-Type", "application/json")
		uh.Logger.ZapLogger.Error("Error to handle request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	user, err := uh.UserService.LoginUser(loginReq)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		respErr := apperrors.MatchError(err)
		w.WriteHeader(respErr.Status)
		_ = json.NewEncoder(w).Encode(respErr)
		return
	}
	tokenPair, err := token.CreateTokenPair(user, uh.UserService.Ur.Chache.Config)
	if err != nil {
		respErr := apperrors.MatchError(err)
		w.WriteHeader(respErr.Status)
		_ = json.NewEncoder(w).Encode(respErr)
		return
	}
	ctx := context.Background()
	exp, err := strconv.Atoi(uh.UserService.Ur.Chache.Config.RefExp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = uh.UserService.Ur.Chache.RedisChahe.Set(ctx, tokenPair.Refresh, user.Id, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		uh.Logger.ZapLogger.Error("Error to save token")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Enable to save token"))
		return
	}
	cookie := http.Cookie{
		Name:     "Refresh",
		Value:    tokenPair.Refresh,
		Expires:  time.Now().Add(time.Duration(exp) * time.Minute),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	w.Header().Set("Authorization", "Bearer "+tokenPair.Access)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
