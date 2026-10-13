package domain

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TOKEN_DURATION = time.Hour
const HMAC_SAMPLE_SECRET = "hmacSampleSecret"

type Login struct {
	Username   string         `db:"username"`
	CustomerId sql.NullString `db:"customer_id"`
	Accounts   sql.NullString `db:"account_numbers"`
	Role       string         `db:"role"`
}

type Claims struct {
	CustomerId string   `json:"customer_id"`
	Accounts   []string `json:"accounts"`
	Username   string   `json:"username"`
	Expiry     int64    `json:"exp"`
	Role       string   `json:"role"`
}

func (c Claims) IsUserRole() bool {
	return c.Role == "user"
}

func (c Claims) IsValidAccountId(accountId string) bool {
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

func (c Claims) IsRequestVerifiedWithTokenClaims(urlParams map[string]string) bool {
	if c.CustomerId != urlParams["customer_id"] {
		return false
	}

	if !c.IsValidAccountId(urlParams["account_id"]) {
		return false
	}

	return true
}

func (l Login) GenerateToken() (string, error) {
	var claims jwt.MapClaims
	if l.Accounts.Valid && l.CustomerId.Valid {
		claims = l.ClaimsForUser()
	} else {
		claims = l.ClaimsForAdmin()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedTokenAsString, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		log.Println("Failed while signing token: " + err.Error())
		return "", errors.New("cannot generate token")
	}

	return signedTokenAsString, nil
}

func (l Login) ClaimsForAccessToken() Claims {
	if l.Accounts.Valid && l.CustomerId.Valid {
		accounts := []string{}
		if l.Accounts.String != "" {
			accounts = strings.Split(l.Accounts.String, ",")
		}
		return Claims{
			CustomerId: l.CustomerId.String,
			Accounts:   accounts,
			Username:   l.Username,
			Expiry:     time.Now().Add(TOKEN_DURATION).Unix(),
			Role:       l.Role,
		}
	}

	return Claims{
		Username: l.Username,
		Expiry:   time.Now().Add(TOKEN_DURATION).Unix(),
		Role:     l.Role,
	}
}

func (l Login) ClaimsForUser() jwt.MapClaims {
	accounts := strings.Split(l.Accounts.String, ",")
	return jwt.MapClaims{
		"customer_id": l.CustomerId.String,
		"role":        l.Role,
		"username":    l.Username,
		"accounts":    accounts,
		"exp":         time.Now().Add(TOKEN_DURATION).Unix(),
	}
}

func (l Login) ClaimsForAdmin() jwt.MapClaims {
	return jwt.MapClaims{
		"role":     l.Role,
		"username": l.Username,
		"exp":      time.Now().Add(TOKEN_DURATION).Unix(),
	}
}
