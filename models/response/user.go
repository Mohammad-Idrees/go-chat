package response

import (
	db "project/db/sqlc"
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone"`
	CreatedAt time.Time `json:"createdAt"`
}

func BuildUserResponse(user *db.User) *UserResponse {
	return &UserResponse{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	}
}

type LoginUserResponse struct {
	SessionId             uuid.UUID     `json:"sessionId"`
	AccessToken           string        `json:"accessToken"`
	AccessTokenExpiresAt  time.Time     `json:"accessTokenExpiresAt"`
	RefreshToken          string        `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time     `json:"refreshTokenExpiresAt"`
	User                  *UserResponse `json:"user"`
}
