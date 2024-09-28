package user

import (
	"context"
	"database/sql"

	"github.com/malytinKonstantin/go-fiber/internal/db"
	"github.com/malytinKonstantin/go-fiber/internal/shared"
)

type UserRepository struct {
	q *db.Queries
}

var _ shared.Repository = (*UserRepository)(nil)

func NewUserRepository(dbConn *sql.DB) *UserRepository {
	return &UserRepository{
		q: db.New(dbConn),
	}
}

func (r *UserRepository) GetUser(ctx context.Context, id int32) (User, error) {
	user, err := r.q.GetUser(ctx, id)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (r *UserRepository) ListUsers(ctx context.Context) ([]User, error) {
	dbUsers, err := r.q.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]User, len(dbUsers))
	for i, u := range dbUsers {
		users[i] = User{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}
	}
	return users, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, name, email string) (User, error) {
	user, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return User{}, err
	}
	return User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int32) error {
	return r.q.DeleteUser(ctx, id)
}
