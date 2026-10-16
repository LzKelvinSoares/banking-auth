package domain

import "github.com/LzKelvinSoares/lzk-go-lib/errs"

type AuthRepository interface {
	FindBy(username, password string) (*Login, *errs.AppError)
	GenerateAndSaveRefreshTokenToStore(authToken AuthToken) (string, *errs.AppError)
	RefreshTokenExist(refreshToken string) *errs.AppError
}
