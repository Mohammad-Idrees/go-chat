package constants

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrNoRows = pgx.ErrNoRows

var ErrPasswordIncorrect = errors.New("password is incorrect")

var ErrTokenExpired = errors.New("token is expired")
var ErrTokenInvalid = errors.New("token is invalid")

var ErrSessionBlocked = errors.New("session is blocked")
var ErrSessionExpired = errors.New("session is expired")
var ErrLoggedOutSession = errors.New("session is logged out")
var ErrIncorrectSessionUser = errors.New("incorrect session user")
var ErrIncorrectSessionToken = errors.New("incorrect session token")

var ErrAccessDenied = errors.New("resource access denied")

var ErrEmptyAuthHeader = errors.New("authorization header not provided")
var ErrInvalidAuthHeader = errors.New("invalid authorization header format")

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)
