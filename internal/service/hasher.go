package service

import (
	"github.com/absagar/go-bcrypt"
)

func HashPassword(password string) string {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(password, salt)
	return hash
}

func CompareToHash(hash, password string) bool {
	return bcrypt.Match(hash, password)
}
