package service

import (
	"context"
	"project/models"
	"project/models/request"
	"project/models/response"
)

type TokenService interface {
	CreateToken(ctx context.Context, req *request.CreateTokenRequest) (string, *models.JWTClaims, error)
	VerifyToken(ctx context.Context, signedToken string) (*models.JWTClaims, error)
	RenewAccessToken(ctx context.Context, req *request.RenewAccessTokenRequest) (*response.RenewAccessTokenResponse, error)
}
