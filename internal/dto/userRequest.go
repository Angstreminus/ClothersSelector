package dto

type RegisterRequest struct {
	Login    string `json:"login"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
