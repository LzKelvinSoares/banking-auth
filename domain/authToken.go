package domain

import (
	"github.com/LzKelvinSoares/lzk-go-lib/errs"
	"github.com/LzKelvinSoares/lzk-go-lib/logger"
	"github.com/golang-jwt/jwt/v5"
)

type AuthToken struct {
	token *jwt.Token
}

func (a AuthToken) NewAccessToken() (string, *errs.AppError) {
	signedTokenAsString, err := a.token.SignedString([]byte(GetHMACSecret()))
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
	signedString, err := token.SignedString([]byte(GetHMACSecret()))
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

func NewAccessTokenFromRefreshToken(refreshToken string) (string, *errs.AppError) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(GetHMACSecret()), nil
	})
	if err != nil {
		return "", errs.NewAuthenticationError("Invalid or Expired Refresh Token")
	}
	r := token.Claims.(*RefreshTokenClaims)
	accessTokenClaims := r.AccessTokenClaims()
	authToken := NewAuthToken(accessTokenClaims)
	return authToken.NewAccessToken()
}
