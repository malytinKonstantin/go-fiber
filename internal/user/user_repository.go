package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/malytinKonstantin/go-fiber/internal/db"
)

type User struct {
	ID           int32  `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	FullName     string `json:"full_name"`
	Bio          string `json:"bio"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type UserRepository struct {
	q *db.Queries
}

func NewUserRepository(dbConn *sql.DB) *UserRepository {
	return &UserRepository{
		q: db.New(dbConn),
	}
}

func (r *UserRepository) GetUser(ctx context.Context, id int32) (User, error) {
	dbUser, err := r.q.GetUser(ctx, id)
	if err != nil {
		return User{}, err
	}
	return convertDbUserToUser(dbUser), nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (User, error) {
	dbUser, err := r.q.GetUserByUsername(ctx, username)
	if err != nil {
		return User{}, err
	}
	return convertDbUserToUser(dbUser), nil
}

func (r *UserRepository) SearchUsers(ctx context.Context, params db.SearchUsersParams) ([]User, error) {
	dbUsers, err := r.q.SearchUsers(ctx, params)
	if err != nil {
		return nil, err
	}
	return convertDbUsersToUsers(dbUsers), nil
}

func (r *UserRepository) CreateUser(ctx context.Context, params db.CreateUserParams) (User, error) {
	dbUser, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return User{}, err
	}
	return convertDbUserToUser(dbUser), nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, params db.UpdateUserParams) (User, error) {
	dbUser, err := r.q.UpdateUser(ctx, params)
	if err != nil {
		return User{}, err
	}
	return convertDbUserToUser(dbUser), nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int32) error {
	return r.q.DeleteUser(ctx, id)
}

func convertDbUserToUser(dbUser db.Users) User {
	return User{
		ID:           dbUser.ID,
		Username:     dbUser.Username,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		FullName:     dbUser.FullName.String,
		Bio:          dbUser.Bio.String,
		CreatedAt:    dbUser.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    dbUser.UpdatedAt.Format(time.RFC3339),
	}
}

func convertDbUsersToUsers(dbUsers []db.Users) []User {
	users := make([]User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = convertDbUserToUser(dbUser)
	}
	return users
}
