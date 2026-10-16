package domain

import (
	"database/sql"

	"github.com/LzKelvinSoares/lzk-go-lib/errs"
	"github.com/LzKelvinSoares/lzk-go-lib/logger"
	"github.com/jmoiron/sqlx"
)

type AuthRepositorySQLServerDB struct {
	client *sqlx.DB
}

func (d AuthRepositorySQLServerDB) FindBy(username, password string) (*Login, *errs.AppError) {
	var login Login
	sqlVerify := `SELECT 
					u.username,
					u.customer_id,
					u.role,
					STRING_AGG(CAST(a.account_id AS VARCHAR(20)), ',') AS account_numbers
				FROM
					users u
				LEFT JOIN accounts a
					ON a.customer_id = u.customer_id
				WHERE
					u.username = @p1 AND u.password = @p2
				GROUP BY 
					u.customer_id,
					u.username,
					u.role`
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

func (a AuthRepositorySQLServerDB) GenerateAndSaveRefreshTokenToStore(authToken AuthToken) (string, *errs.AppError) {
	// generate refresh token
	var appErr *errs.AppError
	var refreshToken string
	if refreshToken, appErr = authToken.NewRefreshToken(); appErr != nil {
		return "", appErr
	}

	// save it
	sqlInsert := "INSERT INTO refresh_token_store (refresh_token) VALUES (@p1)"
	_, err := a.client.Exec(sqlInsert, refreshToken)
	if err != nil {
		logger.Error("Unexpected DB error: " + err.Error())
		return "", errs.NewUnexpectedError("Unexpected DB error")
	}

	return refreshToken, nil
}

func (d AuthRepositorySQLServerDB) RefreshTokenExist(refreshToken string) *errs.AppError {
	sqlSelect := `SELECT 
					refresh_token 
				FROM 
					refresh_token_store 
				WHERE 
					refresh_token = @p1`
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

func NewAuthRepositorySQLServerDB(dbClient *sqlx.DB) AuthRepositorySQLServerDB {
	return AuthRepositorySQLServerDB{
		client: dbClient,
	}
}
