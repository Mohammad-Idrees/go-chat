package service

import (
	"context"
	"project/models/request"
	"project/models/response"
)

type UserService interface {
	CreateUser(ctx context.Context, req *request.CreateUserRequest) (*response.UserResponse, error)
	GetUserByEmail(ctx context.Context, req *request.GetUserByEmailRequest) (*response.UserResponse, error)
	GetUserById(ctx context.Context, req *request.GetUserByIdRequest) (*response.UserResponse, error)
	LoginUser(ctx context.Context, req *request.LoginUserRequest) (*response.LoginUserResponse, error)
	LogoutUser(ctx context.Context, req *request.LogoutUserRequest) error
	GetUsers(ctx context.Context) (*[]response.UserResponse, error)
}
