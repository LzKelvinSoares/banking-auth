package dto

import (
	"banking-auth/domain"

	"github.com/LzKelvinSoares/lzk-go-lib/logger"
	"github.com/golang-jwt/jwt/v5"
)

type RefreshTokenRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r RefreshTokenRequest) AccessTokenValid() error {
	_, err := jwt.Parse(r.AccessToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})

	if err != nil {
		logger.Error("Error while parsing token: " + err.Error())
		return err
	}
	return nil
}
