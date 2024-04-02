package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Angstreminus/ClothersSelector/config"
	authtoken "github.com/Angstreminus/ClothersSelector/internal/auth/token"
	"github.com/Angstreminus/ClothersSelector/internal/repository"
)

type AuthMiddleware struct {
	Cofig          *config.Config
	UserRepository *repository.UserRepository
}

func NewAuthMiddleware(cfg *config.Config, repo *repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		Cofig:          cfg,
		UserRepository: repo,
	}
}

func (am *AuthMiddleware) ValidateToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header != "" {
			bearerToken := strings.Split(header, " ")
			fmt.Println(bearerToken[1])
			if len(bearerToken) == 2 {
				jwtClaims, err := authtoken.DecodeToken(bearerToken[1], am.Cofig.AccSecr)
				if err != nil {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte("Invalid authorization token - does not match with user creds"))
					return
				}
				login := jwtClaims["sub"].(string)
				exist, err := am.UserRepository.UserExists(login)
				if err != nil || !exist {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte("Invalid login"))
					return
				}
			} else {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Invalid authorization token"))
			}
			next(w, r)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	})
}
