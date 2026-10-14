package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const ACCESS_TOKEN_DURATION = time.Hour
const REFRESH_TOKEN_DURATION = time.Hour * 24 * 30

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	CustomerId string   `json:"customer_id"`
	Accounts   []string `json:"accounts"`
	Username   string   `json:"username"`
	Expiry     int64    `json:"exp"`
	Role       string   `json:"role"`
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	TokenType  string   `json:"token_type"`
	CustomerId string   `json:"cid"`
	Accounts   []string `json:"accounts"`
	Username   string   `json:"un"`
	Expiry     int64    `json:"exp"`
	Role       string   `json:"role"`
}

func (c AccessTokenClaims) RefreshTokenClaims() RefreshTokenClaims {
	return RefreshTokenClaims{
		TokenType:  "refresh_token",
		CustomerId: c.CustomerId,
		Accounts:   c.Accounts,
		Username:   c.Username,
		Role:       c.Role,
		Expiry:     time.Now().Add(REFRESH_TOKEN_DURATION).Unix(),
	}
}

func (c AccessTokenClaims) IsUserRole() bool {
	return c.Role == "user"
}

func (c AccessTokenClaims) IsValidAccountId(accountId string) bool {
	if accountId != "" {
		accountFound := false
		for _, a := range c.Accounts {
			if a == accountId {
				accountFound = true
				break
			}
		}
		return accountFound
	}
	return true
}

func (c AccessTokenClaims) IsRequestVerifiedWithTokenClaims(urlParams map[string]string) bool {
	if c.CustomerId != urlParams["customer_id"] {
		return false
	}

	if !c.IsValidAccountId(urlParams["account_id"]) {
		return false
	}

	return true
}

func (c RefreshTokenClaims) AccessTokenClaims() AccessTokenClaims {
	return AccessTokenClaims{
		CustomerId: c.CustomerId,
		Accounts:   c.Accounts,
		Username:   c.Username,
		Role:       c.Role,
		Expiry:     time.Now().Add(ACCESS_TOKEN_DURATION).Unix(),
	}
}
