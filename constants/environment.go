package constants

import (
	"fmt"
	"os"
)

type ServerVariables struct {
	Address string
	Port    string
}

type DBVariables struct {
	Username string
	Password string
	Address  string
	Port     string
	Name     string
}

func (d DBVariables) GetDataSource() string {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		d.Username,
		d.Password,
		d.Address,
		d.Port,
		d.Name)
	return dataSource
}

func GetServerVariables() ServerVariables {

	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")

	return ServerVariables{
		Address: address,
		Port:    port,
	}
}

func GetDBVariables() DBVariables {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbAddress := os.Getenv("DB_ADDRESS")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	return DBVariables{
		Username: dbUser,
		Password: dbPassword,
		Address:  dbAddress,
		Port:     dbPort,
		Name:     dbName,
	}
}
