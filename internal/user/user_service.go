package user

import (
	"context"
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

func (s *UserService) ListUsers(ctx context.Context) ([]User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *UserService) CreateUser(ctx context.Context, name, email string) (User, error) {
	return s.repo.CreateUser(ctx, name, email)
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	return s.repo.DeleteUser(ctx, id)
}
