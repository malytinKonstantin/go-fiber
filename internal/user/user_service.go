package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/malytinKonstantin/go-fiber/internal/auth"
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

func (s *UserService) SearchUsers(ctx context.Context, params SearchUsersParams) ([]User, error) {
	dbParams := db.SearchUsersParams{
		Username:    params.Username,
		Email:       params.Email,
		FullName:    params.FullName,
		Bio:         params.Bio,
		SortBy:      params.SortBy,
		LimitParam:  sql.NullInt32{Int32: params.Limit, Valid: true},
		OffsetParam: sql.NullInt32{Int32: params.Offset, Valid: true},
	}

	if params.CreatedFrom != "" {
		createdFrom, err := time.Parse("2006-01-02", params.CreatedFrom)
		if err != nil {
			return nil, fmt.Errorf("неверный формат даты CreatedFrom: %w", err)
		}
		dbParams.CreatedFrom = &createdFrom
	}

	if params.CreatedTo != "" {
		createdTo, err := time.Parse("2006-01-02", params.CreatedTo)
		if err != nil {
			return nil, fmt.Errorf("неверный формат даты CreatedTo: %w", err)
		}
		dbParams.CreatedTo = &createdTo
	}

	return s.repo.SearchUsers(ctx, dbParams)
}

func (s *UserService) CreateUser(ctx context.Context, dto CreateUserDto) (User, error) {
	hashedPassword, err := HashPassword(dto.Password)
	if err != nil {
		return User{}, err
	}
	dbParams := db.CreateUserParams{
		Username:     dto.Username,
		Email:        dto.Email,
		PasswordHash: hashedPassword,
		FullName:     sql.NullString{String: dto.FullName.String, Valid: dto.FullName.Valid},
		Bio:          sql.NullString{String: dto.Bio.String, Valid: dto.Bio.Valid},
	}
	return s.repo.CreateUser(ctx, dbParams)
}

func (s *UserService) UpdateUser(ctx context.Context, id int32, dto UpdateUserDto) (User, error) {
	dbParams := db.UpdateUserParams{
		ID: id,
	}

	if dto.Username.Valid {
		dbParams.Username = dto.Username.String
	}
	if dto.Email.Valid {
		dbParams.Email = dto.Email.String
	}
	if dto.Password.Valid {
		hashedPassword, err := HashPassword(dto.Password.String)
		if err != nil {
			return User{}, err
		}
		dbParams.PasswordHash = hashedPassword
	}
	if dto.FullName.Valid {
		dbParams.FullName = sql.NullString{String: dto.FullName.String, Valid: true}
	}
	if dto.Bio.Valid {
		dbParams.Bio = sql.NullString{String: dto.Bio.String, Valid: true}
	}

	return s.repo.UpdateUser(ctx, dbParams)
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

func (s *UserService) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	token, err := auth.GenerateToken(auth.User{ID: user.ID})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) ValidateToken(tokenString string) (*auth.Claims, error) {
	return auth.ValidateToken(tokenString)
}

type SearchUsersParams struct {
	Username    string
	Email       string
	FullName    string
	Bio         string
	CreatedFrom string
	CreatedTo   string
	SortBy      string
	Limit       int32
	Offset      int32
}

type CreateUserParams struct {
	Username string
	Email    string
	Password string
	FullName string
	Bio      string
}

type UpdateUserParams struct {
	ID           int32
	Username     string
	Email        string
	PasswordHash string
	FullName     string
	Bio          string
}
