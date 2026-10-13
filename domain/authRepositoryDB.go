package domain

import (
	"database/sql"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	FindBy(username, password string) (*Login, error)
}

type AuthRepositoryDB struct {
	client *sqlx.DB
}

func (d AuthRepositoryDB) FindBy(username, password string) (*Login, error) {
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
			return nil, errors.New("invalid credentials")

		} else {
			log.Println("Error while verifying login from DB: " + err.Error())
			return nil, errors.New("Unexpected DB error")
		}
	}

	return &login, nil
}

func NewAuthRepositoryDB(dbClient *sqlx.DB) AuthRepositoryDB {
	return AuthRepositoryDB{
		client: dbClient,
	}
}
