package domain

import (
	"banking-auth/logger"
	"database/sql"

	"github.com/LzKelvinSoares/lzk-go-lib/errs"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	FindBy(username, password string) (*Login, *errs.AppError)
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

func NewAuthRepositoryDB(dbClient *sqlx.DB) AuthRepositoryDB {
	return AuthRepositoryDB{
		client: dbClient,
	}
}
