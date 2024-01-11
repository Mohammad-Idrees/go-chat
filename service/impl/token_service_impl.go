package service

import (
	"context"
	"project/config"
	"project/constants"
	db "project/db/sqlc"
	"project/logger"
	"project/models"
	"project/models/request"
	"project/models/response"
	"project/service"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenServiceImpl struct {
	jwtSecret           string
	repo                db.Repository
	accessTokenDuration time.Duration
}

func ConfigureTokenService(cfg *config.StartupConfig, repo db.Repository) service.TokenService {
	return &TokenServiceImpl{cfg.Token.JWTSecret, repo, cfg.Token.AccessTokenDuration}
}

func NewJWTClaims(ctx context.Context, req *request.CreateTokenRequest) (*models.JWTClaims, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		logger.Error(ctx, "NewJWTClaims :: failed generating uuid", logger.Field("error", err.Error()))
		return nil, err
	}

	return &models.JWTClaims{
		Id:       tokenId,
		Email:    req.Email,
		IssuedAt: time.Now(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(req.ExpiresIn)),
		},
	}, nil
}

// CreateToken implements service.TokenService.
func (svc *TokenServiceImpl) CreateToken(ctx context.Context, req *request.CreateTokenRequest) (string, *models.JWTClaims, error) {
	jwtClaim, err := NewJWTClaims(ctx, req)
	if err != nil {
		logger.Error(ctx, "CreateToken :: failed creating jwt claims", logger.Field("error", err.Error()))
		return "", nil, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)
	signedToken, err := jwtToken.SignedString([]byte(svc.jwtSecret))
	if err != nil {
		logger.Error(ctx, "CreateToken :: failed signing jwt token", logger.Field("error", err.Error()))
		return "", nil, err
	}

	return signedToken, jwtClaim, nil
}

// VerifyToken implements service.TokenService.
func (svc *TokenServiceImpl) VerifyToken(ctx context.Context, signedToken string) (*models.JWTClaims, error) {
	jwtClaims := &models.JWTClaims{}
	options := []jwt.ParserOption{jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}), jwt.WithExpirationRequired()}
	keyfunc := func(t *jwt.Token) (interface{}, error) {
		return []byte(svc.jwtSecret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(signedToken, jwtClaims, keyfunc, options...)
	if err != nil {
		logger.Error(ctx, "VerifyToken :: failed to parse token with claims", logger.Field("error", err.Error()))
		if strings.Contains(err.Error(), jwt.ErrTokenExpired.Error()) {
			return nil, constants.ErrTokenExpired
		}
		return nil, constants.ErrTokenInvalid
	}

	if !jwtToken.Valid {
		logger.Error(ctx, "VerifyToken :: token is invalid")
		return nil, constants.ErrTokenInvalid
	}

	return jwtClaims, nil
}

// RenewAccessToken implements service.TokenService.
func (svc *TokenServiceImpl) RenewAccessToken(ctx context.Context, req *request.RenewAccessTokenRequest) (*response.RenewAccessTokenResponse, error) {
	refreshClaims, err := svc.VerifyToken(ctx, req.RefreshToken)
	if err != nil {
		logger.Error(ctx, "RenewAccessToken :: failed getting refresh claims", logger.Field("error", err.Error()))
		return nil, err
	}

	session, err := svc.repo.GetSession(ctx, refreshClaims.Id)
	if err != nil {
		logger.Error(ctx, "RenewAccessToken :: failed getting session", logger.Field("error", err.Error()))
		return nil, err
	}

	if session.IsBlocked {
		return nil, constants.ErrSessionBlocked
	}

	if session.Email != refreshClaims.Email {
		return nil, constants.ErrIncorrectSessionUser
	}

	if session.RefreshToken != req.RefreshToken {
		return nil, constants.ErrIncorrectSessionToken
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, constants.ErrSessionExpired
	}

	if session.IsLoggedOut {
		return nil, constants.ErrLoggedOutSession
	}

	accessToken, accessClaims, err := svc.CreateToken(ctx, &request.CreateTokenRequest{Email: refreshClaims.Email, ExpiresIn: svc.accessTokenDuration})
	if err != nil {
		logger.Error(ctx, "RenewAccessToken :: failed to create access token", logger.Field("error", err.Error()))
		return nil, err
	}

	res := &response.RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.ExpiresAt.Time,
	}

	return res, nil

}
