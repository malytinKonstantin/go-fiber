package services

import (
	"context"

	"github.com/malytinKonstantin/sqlc-test/internal/models"
	"github.com/malytinKonstantin/sqlc-test/internal/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, id int32) (models.User, error) {
	return s.repo.GetUser(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *UserService) CreateUser(ctx context.Context, name, email string) (models.User, error) {
	return s.repo.CreateUser(ctx, name, email)
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	return s.repo.DeleteUser(ctx, id)
}
