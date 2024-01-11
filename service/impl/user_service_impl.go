package service

import (
	"context"
	"project/config"
	"project/constants"
	db "project/db/sqlc"
	"project/logger"
	"project/models/request"
	"project/models/response"
	"project/service"
	"project/utils"
	"time"
)

type UserServiceImpl struct {
	repo                 db.Repository
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	tokenSvc             service.TokenService
}

func ConfigureUserService(cfg *config.StartupConfig, repo db.Repository, tokenSvc service.TokenService) service.UserService {
	return &UserServiceImpl{repo, cfg.Token.AccessTokenDuration, cfg.Token.RefreshTokenDuration, tokenSvc}
}

// CreateUser implements service.UserService.
func (svc *UserServiceImpl) CreateUser(ctx context.Context, req *request.CreateUserRequest) (*response.UserResponse, error) {

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Error(ctx, "CreateUser :: failed to hash password", logger.Field("error", err.Error()))
		return nil, err
	}

	arg := &db.CreateUserParams{
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPassword,
		Phone:          req.Phone,
	}

	user, err := svc.repo.CreateUser(ctx, arg)
	if err != nil {
		logger.Error(ctx, "CreateUser :: failed to create user", logger.Field("error", err.Error()))
		return nil, err
	}

	return response.BuildUserResponse(user), nil
}

// GetUser implements service.UserService.
func (svc *UserServiceImpl) GetUserByEmail(ctx context.Context, req *request.GetUserByEmailRequest) (*response.UserResponse, error) {
	user, err := svc.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logger.Error(ctx, "GetUser :: failed to get user", logger.Field("error", err.Error()))
		return nil, err
	}

	return response.BuildUserResponse(user), nil
}

// GetUser implements service.UserService.
func (svc *UserServiceImpl) GetUserById(ctx context.Context, req *request.GetUserByIdRequest) (*response.UserResponse, error) {
	user, err := svc.repo.GetUserById(ctx, req.Id)
	if err != nil {
		logger.Error(ctx, "GetUser :: failed to get user", logger.Field("request", req), logger.Field("error", err.Error()))
		return nil, err
	}

	return response.BuildUserResponse(user), nil
}

// LoginUser implements service.UserService.
func (svc *UserServiceImpl) LoginUser(ctx context.Context, req *request.LoginUserRequest) (*response.LoginUserResponse, error) {
	user, err := svc.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logger.Error(ctx, "LoginUser :: failed to get user", logger.Field("error", err.Error()))
		return nil, err
	}

	err = utils.VerifyPassword(user.HashedPassword, req.Password)
	if err != nil {
		logger.Error(ctx, "LoginUser :: incorrect password", logger.Field("error", err.Error()))
		return nil, constants.ErrPasswordIncorrect
	}

	accessToken, accessClaims, err := svc.tokenSvc.CreateToken(ctx, &request.CreateTokenRequest{Email: req.Email, ExpiresIn: svc.accessTokenDuration})
	if err != nil {
		logger.Error(ctx, "LoginUser :: failed to create access token", logger.Field("error", err.Error()))
		return nil, err
	}

	refreshToken, refreshClaims, err := svc.tokenSvc.CreateToken(ctx, &request.CreateTokenRequest{Email: req.Email, ExpiresIn: svc.refreshTokenDuration})
	if err != nil {
		logger.Error(ctx, "LoginUser :: failed to create refresh token", logger.Field("error", err.Error()))
		return nil, err
	}

	arg := &db.CreateSessionParams{
		ID:           refreshClaims.Id,
		Email:        req.Email,
		UserAgent:    req.UserAgent,
		ClientIp:     req.ClientIp,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshClaims.ExpiresAt.Time,
	}
	session, err := svc.repo.CreateSession(ctx, arg)
	if err != nil {
		logger.Error(ctx, "LoginUser :: failed to create session", logger.Field("error", err.Error()))
		return nil, err
	}

	loginUserResponse := response.LoginUserResponse{
		SessionId:             session.ID, // session.ID = refreshClaims.Id
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
		User:                  response.BuildUserResponse(user),
	}
	return &loginUserResponse, nil
}

// LogoutUser implements service.UserService.
func (svc *UserServiceImpl) LogoutUser(ctx context.Context, req *request.LogoutUserRequest) error {
	refreshClaims, err := svc.tokenSvc.VerifyToken(ctx, req.RefreshToken)
	if err != nil {
		logger.Error(ctx, "LogoutUser :: failed getting refresh claims", logger.Field("error", err.Error()))
		return err
	}

	session, err := svc.repo.GetSession(ctx, refreshClaims.Id)
	if err != nil {
		logger.Error(ctx, "LogoutUser :: failed getting session", logger.Field("error", err.Error()))
		return err
	}

	if session.IsBlocked {
		return constants.ErrSessionBlocked
	}

	if session.Email != refreshClaims.Email {
		return constants.ErrIncorrectSessionUser
	}

	if session.RefreshToken != req.RefreshToken {
		return constants.ErrIncorrectSessionToken
	}

	if time.Now().After(session.ExpiresAt) {
		return constants.ErrSessionExpired
	}

	if session.IsLoggedOut {
		return constants.ErrLoggedOutSession
	}

	_, err = svc.repo.UpdateSession(ctx, &db.UpdateSessionParams{IsLoggedOut: utils.BoolPtr(true), ID: session.ID})
	if err != nil {
		logger.Error(ctx, "LogoutUser :: failed updating session", logger.Field("error", err.Error()))
		return err
	}

	return nil
}

// GetUsers implements service.UserService.
func (svc *UserServiceImpl) GetUsers(ctx context.Context) (*[]response.UserResponse, error) {
	users, err := svc.repo.GetUsers(ctx)
	if err != nil {
		logger.Error(ctx, "GetUsers :: failed to get user", logger.Field("error", err.Error()))
		return nil, err
	}

	userResp := make([]response.UserResponse, 0)
	for _, user := range users {
		userResp = append(userResp, *response.BuildUserResponse(user))
	}

	return &userResp, nil
}
