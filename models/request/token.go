package request

import "time"

type CreateTokenRequest struct {
	Email     string
	Type      string
	ExpiresIn time.Duration
}

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
