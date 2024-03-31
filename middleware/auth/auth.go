package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Angstreminus/ClothersSelector/config"
	"github.com/Angstreminus/ClothersSelector/internal/apperrors"
	"github.com/Angstreminus/ClothersSelector/internal/dto"
	"github.com/Angstreminus/ClothersSelector/internal/repository"
	"github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
	Cofig          *config.Config
	UserRepository *repository.UserRepository
}

func (am *AuthMiddleware) ValidateToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header != "" {
			bearerToken := strings.Split(header, " ")
			if len(bearerToken) == 2 {
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, &apperrors.TokenError{
							Message: "Token error",
						}
					}
					return []byte(am.Cofig.AccSecr), nil
				})
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(err.Error())
					return
				}
				if token.Valid {
					claims, ok := token.Claims.(jwt.MapClaims)
					if !ok && !token.Valid {
						w.WriteHeader(http.StatusForbidden)
					}
					var usrSign dto.UserSignature
					usrSign.Login = claims["login"].(string)
					exist, err := am.UserRepository.UserExists(usrSign.Login)
					if err != nil {
						w.WriteHeader(http.StatusForbidden)
						w.Write([]byte("Invalid authorization token - does not match with user creds"))
						return
					}
					if !exist {
						w.WriteHeader(http.StatusForbidden)
						w.Write([]byte("Invalid authorization token"))
					}
					next(w, r)
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})
}
