package authtoken

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Angstreminus/ClothersSelector/config"
	"github.com/Angstreminus/ClothersSelector/internal/apperrors"
	"github.com/Angstreminus/ClothersSelector/internal/entity"
	"github.com/golang-jwt/jwt/v4"
)

type TokenPair struct {
	Access  string
	Refresh string
}

func DecodeToken(tokenString, secret string) (jwt.MapClaims, apperrors.AppError) {
	secretBytes := []byte(secret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &apperrors.TokenError{
				Message: "Wrong method",
			}
		}
		return secretBytes, nil
	})
	if err != nil {
		fmt.Println("Error to parse bytes")
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		fmt.Println("Wrong token")
		return nil, &apperrors.TokenError{
			Message: "Wrong token",
		}
	}
}

func CreateToken(user entity.User, expTime, secret string) (string, apperrors.AppError) {
	tokenExpTime, err := strconv.Atoi(expTime)
	if err != nil {
		return "", &apperrors.TokenError{
			Message: err.Error(),
		}
	}
	payload := jwt.MapClaims{
		"sub": user.Login,
		"exp": time.Now().Add(time.Minute * time.Duration(tokenExpTime)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", &apperrors.TokenError{
			Message: err.Error(),
		}
	}
	return t, err
}

func CreateTokenPair(user entity.User, cfg config.Config) (TokenPair, apperrors.AppError) {
	accToken, err := CreateToken(user, cfg.AccExp, cfg.AccSecr)
	if err != nil {
		fmt.Println("Create token error")
		return TokenPair{}, err
	}

	refToken, err := CreateToken(user, cfg.RefExp, cfg.RefSecr)
	if err != nil {
		fmt.Println("Create token error")
		return TokenPair{}, err
	}
	return TokenPair{
		Access:  accToken,
		Refresh: refToken,
	}, nil
}
