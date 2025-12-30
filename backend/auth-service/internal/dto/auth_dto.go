package dto

type RegisterRequest struct {
Email     string `json:"email" validate:"required,email"`
Username  string `json:"username" validate:"required,min=3,max=30"`
Password  string `json:"password" validate:"required,min=6"`
FirstName string `json:"firstName" validate:"required"`
LastName  string `json:"lastName" validate:"required"`
}

type LoginRequest struct {
Email    string `json:"email" validate:"required,email"`
Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
RefreshToken string `json:"refreshToken" validate:"required"`
}

type AuthResponse struct {
AccessToken  string      `json:"accessToken"`
RefreshToken string      `json:"refreshToken"`
User         interface{} `json:"user"`
}

type ErrorResponse struct {
Error   string `json:"error"`
Message string `json:"message"`
}

type SuccessResponse struct {
Message string      `json:"message"`
Data    interface{} `json:"data,omitempty"`
}
