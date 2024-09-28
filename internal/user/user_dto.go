package user

import (
	"github.com/malytinKonstantin/go-fiber/internal/shared"
)

type CreateUserParams struct {
	Username string            `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email    string            `json:"email" validate:"required,email,max=100"`
	Password string            `json:"password" validate:"required,min=8,max=20,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=0123456789,containsany=!@#$%^&*()"`
	FullName shared.NullString `json:"full_name" validate:"omitempty,max=100"`
	Bio      shared.NullString `json:"bio" validate:"omitempty,max=500"`
}

type SignInInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignInOutput struct {
	Token string `json:"token"`
}
