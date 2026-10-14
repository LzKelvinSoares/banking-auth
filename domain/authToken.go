package domain

import (
	"banking-auth/logger"

	"github.com/LzKelvinSoares/lzk-go-lib/errs"
	"github.com/golang-jwt/jwt/v5"
)

type AuthToken struct {
	token *jwt.Token
}

func (a AuthToken) NewAccessToken() (string, *errs.AppError) {
	signedTokenAsString, err := a.token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		logger.Error("Failed while signing token: " + err.Error())
		return "", errs.NewUnexpectedError("cannot generate token")
	}

	return signedTokenAsString, nil
}

func (a AuthToken) NewRefreshToken() (string, *errs.AppError) {
	c := a.token.Claims.(AccessTokenClaims)
	refreshClaims := c.RefreshTokenClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedString, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		logger.Error("Failed while refresh token: " + err.Error())
		return "", errs.NewUnexpectedError("cannot generate refresh token")
	}

	return signedString, nil
}

func NewAuthToken(claims AccessTokenClaims) AuthToken {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return AuthToken{
		token: token,
	}
}
