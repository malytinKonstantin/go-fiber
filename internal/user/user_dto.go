package user

import (
	"github.com/malytinKonstantin/go-fiber/internal/shared"
)

type CreateUserDto struct {
	Username string            `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email    string            `json:"email" validate:"required,email,max=100"`
	Password string            `json:"password" validate:"required,min=8,max=20,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=0123456789,containsany=!@#$%^&*()"`
	FullName shared.NullString `json:"full_name" validate:"omitempty,max=100"`
	Bio      shared.NullString `json:"bio" validate:"omitempty,max=500"`
}

type UpdateUserDto struct {
	Username shared.NullString `json:"username" validate:"omitempty,min=3,max=50,alphanum"`
	Email    shared.NullString `json:"email" validate:"omitempty,email,max=100"`
	Password shared.NullString `json:"password" validate:"omitempty,min=8,max=20,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=0123456789,containsany=!@#$%^&*()"`
	FullName shared.NullString `json:"full_name" validate:"omitempty,max=100"`
	Bio      shared.NullString `json:"bio" validate:"omitempty,max=500"`
}

type SignInDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignInOutput struct {
	Token string `json:"token"`
}

type ListUsersQuery struct {
	Username    string `query:"username"`
	Email       string `query:"email"`
	FullName    string `query:"full_name"`
	Bio         string `query:"bio"`
	CreatedFrom string `query:"created_from"`
	CreatedTo   string `query:"created_to"`
	SortBy      string `query:"sort_by"`
	Limit       int32  `query:"limit"`
	Offset      int32  `query:"offset"`
}
