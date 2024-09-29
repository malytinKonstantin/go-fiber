package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/malytinKonstantin/go-fiber/internal/auth"
	"github.com/malytinKonstantin/go-fiber/internal/db"
	"golang.org/x/crypto/bcrypt"
)

const (
	dateFormat            = "2006-01-02"
	invalidDateFormatErr  = "invalid date format"
	invalidCredentialsErr = "invalid credentials"
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, id int32) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}
	return s.repo.GetUser(ctx, id)
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}
	return s.repo.GetUserByUsername(ctx, username)
}

func (s *UserService) SearchUsers(ctx context.Context, params SearchUsersParams) ([]User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	dbParams := db.SearchUsersParams{
		Username:    params.Username,
		Email:       params.Email,
		FullName:    params.FullName,
		Bio:         params.Bio,
		SortBy:      params.SortBy,
		LimitParam:  sql.NullInt32{Int32: params.Limit, Valid: params.Limit > 0},
		OffsetParam: sql.NullInt32{Int32: params.Offset, Valid: params.Offset >= 0},
	}

	if params.CreatedFrom != "" || params.CreatedTo != "" {
		if err := s.parseAndSetDates(&dbParams, params.CreatedFrom, params.CreatedTo); err != nil {
			return nil, err
		}
	}

	return s.repo.SearchUsers(ctx, dbParams)
}

func (s *UserService) parseAndSetDates(dbParams *db.SearchUsersParams, createdFrom, createdTo string) error {
	if createdFrom != "" {
		if t, err := time.Parse(dateFormat, createdFrom); err != nil {
			return errors.New(invalidDateFormatErr)
		} else {
			tPtr := &t
			dbParams.CreatedFrom = &tPtr
		}
	}

	if createdTo != "" {
		if t, err := time.Parse(dateFormat, createdTo); err != nil {
			return errors.New(invalidDateFormatErr)
		} else {
			tPtr := &t
			dbParams.CreatedTo = &tPtr
		}
	}

	return nil
}

func (s *UserService) CreateUser(ctx context.Context, dto CreateUserDto) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}

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
	if err := ctx.Err(); err != nil {
		return User{}, err
	}

	dbParams := db.UpdateUserParams{ID: id}

	if err := s.setUpdateParams(&dbParams, dto); err != nil {
		return User{}, err
	}

	return s.repo.UpdateUser(ctx, dbParams)
}

func (s *UserService) setUpdateParams(dbParams *db.UpdateUserParams, dto UpdateUserDto) error {
	if dto.Username.Valid {
		dbParams.Username = dto.Username.String
	}
	if dto.Email.Valid {
		dbParams.Email = dto.Email.String
	}
	if dto.Password.Valid {
		hashedPassword, err := HashPassword(dto.Password.String)
		if err != nil {
			return err
		}
		dbParams.PasswordHash = hashedPassword
	}
	if dto.FullName.Valid {
		dbParams.FullName = sql.NullString{String: dto.FullName.String, Valid: true}
	}
	if dto.Bio.Valid {
		dbParams.Bio = sql.NullString{String: dto.Bio.String, Valid: true}
	}
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	if err := ctx.Err(); err != nil {
		return err
	}
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
	if err := ctx.Err(); err != nil {
		return "", err
	}

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New(invalidCredentialsErr)
		}
		return "", err
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New(invalidCredentialsErr)
	}

	return auth.GenerateToken(auth.User{ID: user.ID})
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
