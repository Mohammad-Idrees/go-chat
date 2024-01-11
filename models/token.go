package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TODO: add type to claims to differentiate between refresh token and refresh token
type JWTClaims struct {
	Id       uuid.UUID
	Email    string
	IssuedAt time.Time
	jwt.RegisteredClaims
}
