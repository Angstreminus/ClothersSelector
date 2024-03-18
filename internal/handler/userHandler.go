package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Angstreminus/ClothersSelector/internal/apperrors"
	"github.com/Angstreminus/ClothersSelector/internal/auth/token"
	"github.com/Angstreminus/ClothersSelector/internal/dto"
	"github.com/Angstreminus/ClothersSelector/internal/service"
)

type Cookie struct {
	Refresh string `json:"login" validate:"required"`
}

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(usrServ *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: usrServ,
	}
}

func (uh *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var toRegistrate *dto.ReisterRequest
	if err := json.NewDecoder(r.Body).Decode(toRegistrate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	user, err := uh.UserService.RegiterUser(toRegistrate)
	if err != nil {
		respErr := apperrors.MatchError(err)
		w.WriteHeader(respErr.Status)
		_ = json.NewEncoder(w).Encode(respErr)
	}
	tokenPair, err := token.CreateTokenPair(user, uh.UserService.Ur.Chache.Config)
	if err != nil {
		respErr := apperrors.MatchError(err)
		w.WriteHeader(respErr.Status)
		_ = json.NewEncoder(w).Encode(respErr)
	}
	ctx := context.Background()
	exp, err := strconv.Atoi(uh.UserService.Ur.Chache.Config.RefExp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	uh.UserService.Ur.Chache.RedisChahe.Set(ctx, tokenPair.Refresh, user.Id.String(), time.Minute*time.Duration(exp))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization:", "Bearer "+tokenPair.Access)
	cookie := http.Cookie{
		Name:  "Refresh",
		Value: tokenPair.Refresh,
	}
	http.SetCookie(w, &cookie)
	w.Header().Set("Authorization", "bearer access_token")
	if err = json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
