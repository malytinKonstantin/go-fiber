package user

import (
	"context"
	"database/sql"

	"github.com/malytinKonstantin/go-fiber/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, id int32) (User, error) {
	return s.repo.GetUser(ctx, id)
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (User, error) {
	return s.repo.GetUserByUsername(ctx, username)
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int32) ([]User, error) {
	return s.repo.ListUsers(ctx, limit, offset)
}

func (s *UserService) CreateUser(ctx context.Context, params CreateUserParams) (User, error) {
	hashedPassword, err := HashPassword(params.Password)
	if err != nil {
		return User{}, err
	}
	dbParams := db.CreateUserParams{
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: hashedPassword,
		FullName:     sql.NullString{String: params.FullName.String, Valid: params.FullName.Valid},
		Bio:          sql.NullString{String: params.Bio.String, Valid: params.Bio.Valid},
	}
	return s.repo.CreateUser(ctx, dbParams)
}

func (s *UserService) UpdateUser(ctx context.Context, params db.UpdateUserParams) (User, error) {
	return s.repo.UpdateUser(ctx, params)
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	return s.repo.DeleteUser(ctx, id)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
