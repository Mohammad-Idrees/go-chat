package response

import (
	"time"
)

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"accessToke"`
	AccessTokenExpiresAt time.Time `json:"expiesAt"`
}
