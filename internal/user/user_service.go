package user

import (
	"context"

	"github.com/malytinKonstantin/go-fiber/internal/db"
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

func (s *UserService) CreateUser(ctx context.Context, params db.CreateUserParams) (User, error) {
	return s.repo.CreateUser(ctx, params)
}

func (s *UserService) UpdateUser(ctx context.Context, params db.UpdateUserParams) (User, error) {
	return s.repo.UpdateUser(ctx, params)
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	return s.repo.DeleteUser(ctx, id)
}
