package user

import "github.com/malytinKonstantin/go-fiber/internal/shared"

type CreateUserParams struct {
	Username string            `json:"username"`
	Email    string            `json:"email"`
	Password string            `json:"password"`
	FullName shared.NullString `json:"full_name"`
	Bio      shared.NullString `json:"bio"`
}
