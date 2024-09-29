package user

import (
	"github.com/malytinKonstantin/go-fiber/internal/shared"
)

// CreateUserDto represents the data for creating a user
// swagger:model
type CreateUserDto struct {
	// Username of the user
	// required: true
	// min: 3
	// max: 50
	// example: johndoe
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`

	// Email of the user
	// required: true
	// max: 100
	// example: john@example.com
	Email string `json:"email" validate:"required,email,max=100"`

	// Password of the user
	// required: true
	// min: 8
	// max: 20
	// example: P@ssw0rd!
	Password string `json:"password" validate:"required,min=8,max=20,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=0123456789,containsany=!@#$%^&*()"`

	// Full name of the user
	// max: 100
	// example: John Doe
	FullName shared.NullString `json:"full_name" validate:"omitempty,max=100"`

	// Biography of the user
	// max: 500
	// example: Software developer with 5 years of experience
	Bio shared.NullString `json:"bio" validate:"omitempty,max=500"`
}

// UpdateUserDto represents the data for updating a user
// swagger:model
type UpdateUserDto struct {
	// Username of the user
	// min: 3
	// max: 50
	// example: johndoe
	Username shared.NullString `json:"username" validate:"omitempty,min=3,max=50,alphanum"`

	// Email of the user
	// max: 100
	// example: john@example.com
	Email shared.NullString `json:"email" validate:"omitempty,email,max=100"`

	// Password of the user
	// min: 8
	// max: 20
	// example: NewP@ssw0rd!
	Password shared.NullString `json:"password" validate:"omitempty,min=8,max=20,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=0123456789,containsany=!@#$%^&*()"`

	// Full name of the user
	// max: 100
	// example: John Doe
	FullName shared.NullString `json:"full_name" validate:"omitempty,max=100"`

	// Biography of the user
	// max: 500
	// example: Experienced software developer and team lead
	Bio shared.NullString `json:"bio" validate:"omitempty,max=500"`
}

// SignInDto represents the data for user sign-in
// swagger:model
type SignInDto struct {
	// Username of the user
	// required: true
	// example: johndoe
	Username string `json:"username" validate:"required"`

	// Password of the user
	// required: true
	// example: P@ssw0rd!
	Password string `json:"password" validate:"required"`
}

// SignInOutput represents the result of a successful sign-in
// swagger:model
type SignInOutput struct {
	// JWT token
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	Token string `json:"token"`
}

// ListUsersQuery represents the query parameters for listing users
// swagger:model
type ListUsersQuery struct {
	// Username for filtering
	// example: john
	Username string `query:"username"`

	// Email for filtering
	// example: john@example.com
	Email string `query:"email"`

	// Full name for filtering
	// example: John Doe
	FullName string `query:"full_name"`

	// Biography for filtering
	// example: developer
	Bio string `query:"bio"`

	// Creation date (from) in YYYY-MM-DD format
	// example: 2023-01-01
	CreatedFrom string `query:"created_from"`

	// Creation date (to) in YYYY-MM-DD format
	// example: 2023-12-31
	CreatedTo string `query:"created_to"`

	// Field for sorting
	// example: username_asc
	SortBy string `query:"sort_by"`

	// Limit of returned records
	// example: 10
	Limit int32 `query:"limit"`

	// Offset for pagination
	// example: 0
	Offset int32 `query:"offset"`
}

// ErrorResponse represents the structure of an error response
// swagger:model
type ErrorResponse struct {
	// Error message
	// example: Invalid input
	Error string `json:"error"`
}

// SuccessResponse represents the structure of a successful response
// swagger:model
type SuccessResponse struct {
	// Success message
	// example: Operation completed successfully
	Message string `json:"message"`
}
