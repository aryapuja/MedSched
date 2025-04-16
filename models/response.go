package models

type LoginResponse struct {
	Token    string `json:"token"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
