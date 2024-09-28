package repositories

import (
	"context"
	"database/sql"

	"github.com/malytinKonstantin/sqlc-test/internal/db"
	"github.com/malytinKonstantin/sqlc-test/internal/models"
)

type UserRepository struct {
	q *db.Queries
}

func NewUserRepository(dbConn *sql.DB) *UserRepository {
	return &UserRepository{
		q: db.New(dbConn),
	}
}

func (r *UserRepository) GetUser(ctx context.Context, id int32) (models.User, error) {
	user, err := r.q.GetUser(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return models.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (r *UserRepository) ListUsers(ctx context.Context) ([]models.User, error) {
	dbUsers, err := r.q.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]models.User, len(dbUsers))
	for i, u := range dbUsers {
		users[i] = models.User{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}
	}
	return users, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, name, email string) (models.User, error) {
	user, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return models.User{}, err
	}
	return models.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int32) error {
	return r.q.DeleteUser(ctx, id)
}
