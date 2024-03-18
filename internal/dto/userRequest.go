package dto

type ReisterRequest struct {
	Login    string `json:"login" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
