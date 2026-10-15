package app

import (
	"banking-auth/constants"
	"banking-auth/domain"
	"banking-auth/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func sanityCheck() {
	godotenv.Load()
	if os.Getenv("SERVER_ADDRESS") == "" ||
		os.Getenv("SERVER_PORT") == "" {
		log.Fatal("Environment variables not defined!")
	}
}

func Start() {
	sanityCheck()

	router := mux.NewRouter()
	dbClient := getDBClient()

	//define handlers
	ah := AuthHandler{
		service: service.NewDefaultAuthService(domain.NewAuthRepositoryDB(dbClient), domain.GetRolePermissions()),
	}

	// define routes
	router.HandleFunc(constants.LOGIN, ah.Login).Methods(http.MethodPost)
	router.HandleFunc(constants.VERIFY, ah.Verify).Methods(http.MethodGet)
	router.HandleFunc(constants.REFRESH, ah.Refresh).Methods(http.MethodPost)

	// starting server
	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router))

}

func getDBClient() *sqlx.DB {
	dbVariables := constants.GetDBVariables()

	client, err := sqlx.Open("mysql", dbVariables.GetDataSource())
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}
