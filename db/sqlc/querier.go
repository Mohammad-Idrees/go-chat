// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateChannel(ctx context.Context, name string) (*Channel, error)
	CreateMembership(ctx context.Context, arg *CreateMembershipParams) (*Membership, error)
	CreateSession(ctx context.Context, arg *CreateSessionParams) (*Session, error)
	CreateUser(ctx context.Context, arg *CreateUserParams) (*User, error)
	GetChannelById(ctx context.Context, id int64) (*Channel, error)
	GetChannels(ctx context.Context) ([]*Channel, error)
	GetMemberships(ctx context.Context) ([]*Membership, error)
	GetMembershipsByChannelId(ctx context.Context, channelID int64) ([]*Membership, error)
	GetMembershipsByUserId(ctx context.Context, userID int64) ([]*Membership, error)
	GetSession(ctx context.Context, id uuid.UUID) (*Session, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserById(ctx context.Context, id int64) (*User, error)
	GetUsers(ctx context.Context) ([]*User, error)
	UpdateSession(ctx context.Context, arg *UpdateSessionParams) (*Session, error)
}

var _ Querier = (*Queries)(nil)
