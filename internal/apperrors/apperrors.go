package apperrors

import (
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
	Status  int    `json:"-"`
}

type InvalidDataErr struct {
	Message string
}

func (error *InvalidDataErr) Error() string {
	return error.Message
}

type TokenError struct {
	Message string
}

func (error *TokenError) Error() string {
	return error.Message
}

type HashError struct {
	Message string
}

func (error *HashError) Error() string {
	return error.Message
}

type DBoperationErr struct {
	Message string
}

func (error *DBoperationErr) Error() string {
	return error.Message
}

type AuthError struct {
	Message string
}

func (error *AuthError) Error() string {
	return error.Message
}

type AppError interface {
	Error() string
}

func MatchError(appErr AppError) *ResponseError {
	switch ae := appErr.(type) {
	case *DBoperationErr, *TokenError, *HashError:
		return &ResponseError{
			Message: ae.Error(),
			Status:  http.StatusInternalServerError,
		}
	case *InvalidDataErr:
		return &ResponseError{
			Message: ae.Message,
			Status:  http.StatusBadRequest,
		}
	case *AuthError:
		return &ResponseError{
			Message: ae.Message,
			Status:  http.StatusUnauthorized,
		}
	}
	return nil
}
