package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Angstreminus/ClothersSelector/config"
	"github.com/Angstreminus/ClothersSelector/internal/apperrors"
	authtoken "github.com/Angstreminus/ClothersSelector/internal/auth/token"
	"github.com/Angstreminus/ClothersSelector/internal/dto"
	"github.com/Angstreminus/ClothersSelector/internal/service"
	"github.com/Angstreminus/ClothersSelector/logger"
	"github.com/golang-jwt/jwt/v4"
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
	tokenPair, err := authtoken.CreateTokenPair(user, *uh.Config)
	fmt.Println(tokenPair.Access)
	fmt.Println(tokenPair.Refresh)
	if err != nil {
		fmt.Println("Hanler error after creating pair")
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
		uh.Logger.ZapLogger.Error("Error to handle request")
		w.Header().Set("Content-Type", "application/json")
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
	tokenPair, err := authtoken.CreateTokenPair(*user, *uh.UserService.Ur.Chache.Config)

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

func (uh UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")
	if header != "" {
		bearerToken := strings.Split(header, " ")
		if len(bearerToken) == 2 {
			jwtClaims, err := authtoken.DecodeToken(bearerToken[1], uh.Config.AccSecr)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Invalid authorization token - does not match with user creds"))
				return
			}
			login := jwtClaims["sub"].(string)
			exist, err := uh.UserService.UserExists(login)
			if err != nil || !exist {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Invalid login"))
				return
			}
			jwtClaims["exp"] = time.Now().Unix()
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
			// TODO: End shutdown tokens
			_, err = token.SignedString([]byte(uh.Config.AccSecr))
			if err != nil {
				if err != nil {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte("Invalid authorization token - does not match with user creds"))
					return
				}
			}
			cookie, err := r.Cookie("Refresh")
			if err != nil {
				if err != http.ErrNoCookie {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte(err.Error()))
					return
				}
			}
			cookie.Expires = time.Now()
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Invalid authorization token"))
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func (uh *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	return
}
