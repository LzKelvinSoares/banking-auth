package domain

import (
	"database/sql"
	"strings"
	"time"
)

const HMAC_SAMPLE_SECRET = "hmacSampleSecret"

type Login struct {
	Username   string         `db:"username"`
	CustomerId sql.NullString `db:"customer_id"`
	Accounts   sql.NullString `db:"account_numbers"`
	Role       string         `db:"role"`
}

func (l Login) ClaimsForAccessToken() AccessTokenClaims {
	if l.Accounts.Valid && l.CustomerId.Valid {
		return l.ClaimsForUser()
	}

	return l.ClaimsForAdmin()
}

func (l Login) ClaimsForUser() AccessTokenClaims {
	accounts := strings.Split(l.Accounts.String, ",")
	return AccessTokenClaims{
		CustomerId: l.CustomerId.String,
		Role:       l.Role,
		Username:   l.Username,
		Accounts:   accounts,
		Expiry:     time.Now().Add(ACCESS_TOKEN_DURATION).Unix(),
	}
}

func (l Login) ClaimsForAdmin() AccessTokenClaims {
	return AccessTokenClaims{
		Role:     l.Role,
		Username: l.Username,
		Expiry:   time.Now().Add(ACCESS_TOKEN_DURATION).Unix(),
	}
}
