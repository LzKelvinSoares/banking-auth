package domain

import (
	"banking-auth/logger"
	"database/sql"

	"github.com/LzKelvinSoares/lzk-go-lib/errs"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	FindBy(username, password string) (*Login, *errs.AppError)
	GenerateAndSaveRefreshTokenToStore(authToken AuthToken) (string, *errs.AppError)
	RefreshTokenExist(refreshToken string) *errs.AppError
}

type AuthRepositoryDB struct {
	client *sqlx.DB
}

func (d AuthRepositoryDB) FindBy(username, password string) (*Login, *errs.AppError) {
	var login Login
	sqlVerify := `SELECT 
					u.username,
					u.customer_id,
					u.role,
					GROUP_CONCAT(a.account_id) as account_numbers 
				FROM
					users u
				LEFT JOIN accounts a
					ON a.customer_id = u.customer_id
				WHERE
					u.username = ? AND password = ?
				GROUP BY 
					u.customer_id,
					u.username`
	err := d.client.Get(&login, sqlVerify, username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewAuthenticationError("invalid credentials")

		} else {
			logger.Error("Error while verifying login from DB: " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected DB error")
		}
	}

	return &login, nil
}

func (a AuthRepositoryDB) GenerateAndSaveRefreshTokenToStore(authToken AuthToken) (string, *errs.AppError) {
	// generate refresh token
	var appErr *errs.AppError
	var refreshToken string
	if refreshToken, appErr = authToken.NewRefreshToken(); appErr != nil {
		return "", appErr
	}

	// save it
	sqlInsert := "insert into refresh_token_store (refresh_token) values (?)"
	_, err := a.client.Exec(sqlInsert, refreshToken)
	if err != nil {
		logger.Error("Unexpected DB error: " + err.Error())
		return "", errs.NewUnexpectedError("Unexpected DB error")
	}

	return refreshToken, nil
}

func (d AuthRepositoryDB) RefreshTokenExist(refreshToken string) *errs.AppError {
	sqlSelect := "select refresh_token from refresh_token_store where refresh_token = ?"
	var token string
	if err := d.client.Get(&token, sqlSelect, refreshToken); err != nil {
		if err == sql.ErrNoRows {
			return errs.NewAuthenticationError("Refresh Token not registered")
		} else {
			logger.Error("Unexpected DB error: " + err.Error())
			return errs.NewUnexpectedError("Unexpected DB error")
		}
	}
	return nil
}

func NewAuthRepositoryDB(dbClient *sqlx.DB) AuthRepositoryDB {
	return AuthRepositoryDB{
		client: dbClient,
	}
}
