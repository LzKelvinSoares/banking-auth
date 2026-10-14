package domain

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/LzKelvinSoares/lzk-go-lib/errs"
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

func (l Login) GenerateToken() (string, *errs.AppError) {
	var claims Claims
	if l.Accounts.Valid && l.CustomerId.Valid {
		claims = l.ClaimsForUser()
	} else {
		claims = l.ClaimsForAdmin()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedTokenAsString, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		log.Println("Failed while signing token: " + err.Error())
		return "", errs.NewUnexpectedError("cannot generate token")
	}

	return signedTokenAsString, nil
}

func (l Login) ClaimsForAccessToken() Claims {
	if l.Accounts.Valid && l.CustomerId.Valid {
		return l.ClaimsForUser()
	}

	return l.ClaimsForAdmin()
}

func (l Login) ClaimsForUser() Claims {
	accounts := strings.Split(l.Accounts.String, ",")
	return Claims{
		CustomerId: l.CustomerId.String,
		Role:       l.Role,
		Username:   l.Username,
		Accounts:   accounts,
		Expiry:     time.Now().Add(TOKEN_DURATION).Unix(),
	}
}

func (l Login) ClaimsForAdmin() Claims {
	return Claims{
		Role:     l.Role,
		Username: l.Username,
		Expiry:   time.Now().Add(TOKEN_DURATION).Unix(),
	}
}
