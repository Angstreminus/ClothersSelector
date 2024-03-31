package token

import (
	"strconv"
	"time"

	"github.com/Angstreminus/ClothersSelector/config"
	"github.com/Angstreminus/ClothersSelector/internal/apperrors"
	"github.com/Angstreminus/ClothersSelector/internal/dto"
	"github.com/Angstreminus/ClothersSelector/internal/entity"
	"github.com/golang-jwt/jwt/v4"
)

type CustomToken struct {
	UUID  string `json:"uuid"`
	Login string `json:"login"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	Access  string
	Refresh string
}

func CreateToken(user *entity.User, expTime, secret string) (string, apperrors.AppError) {
	tokenExpTime, err := strconv.Atoi(expTime)
	if err != nil {
		return "", &apperrors.TokenError{
			Message: err.Error(),
		}
	}
	claims := &CustomToken{
		Login: user.Login,
		UUID:  user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Minute * time.Duration(tokenExpTime)),
			},
		},
	}
	templ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := templ.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, err
}

func IsAuthorized(token string, cfg *config.Config) (bool, apperrors.AppError) {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &apperrors.AuthError{
				Message: "INVALID SIGNING METHOD",
			}
		}
		return 0, nil
	})
	if err != nil {
		return false, &apperrors.AuthError{
			Message: err.Error(),
		}
	}
	return true, nil
}

func ExtractFromToken(reqtoken string, cfg *config.Config) (*dto.UserSignature, error) {
	token, err := jwt.Parse(reqtoken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &apperrors.AuthError{
				Message: "INVALID SIGNING METHOD",
			}
		}
		return []byte(cfg.AccSecr), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, &apperrors.AuthError{
			Message: "INVALID TOKEN",
		}
	}
	var usrSign dto.UserSignature
	usrSign.Login = claims["login"].(string)
	usrSign.Password = claims["ID"].(string)
	return &usrSign, nil
}

func CreateTokenPair(user *entity.User, cfg *config.Config) (*TokenPair, apperrors.AppError) {
	accToken, err := CreateToken(user, cfg.AccExp, cfg.AccSecr)
	if err != nil {
		return nil, err
	}

	refToken, err := CreateToken(user, cfg.RefExp, cfg.RefSecr)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		Access:  accToken,
		Refresh: refToken,
	}, nil
}
