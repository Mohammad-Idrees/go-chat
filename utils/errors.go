package utils

import (
	"errors"
	"net/http"
	"project/constants"

	"github.com/jackc/pgx/v5/pgconn"
)

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

func GetHTTPStatusCode(err error) int {
	switch err {
	case constants.ErrPasswordIncorrect, constants.ErrTokenExpired, constants.ErrTokenInvalid, constants.ErrSessionBlocked, constants.ErrSessionExpired,
		constants.ErrLoggedOutSession, constants.ErrIncorrectSessionUser, constants.ErrIncorrectSessionToken, constants.ErrAccessDenied, constants.ErrEmptyAuthHeader,
		constants.ErrInvalidAuthHeader:
		return http.StatusUnauthorized
	case constants.ErrNoRows:
		return http.StatusNotFound
	default:
		errCode := ErrorCode(err)
		if errCode == constants.ForeignKeyViolation || errCode == constants.UniqueViolation {
			return http.StatusUnprocessableEntity
		}
		return http.StatusInternalServerError
	}
}
