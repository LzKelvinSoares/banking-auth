package app

import (
	"banking-auth/constants"
	"banking-auth/domain"
	"banking-auth/service"
	"log"
	"net/http"
	"time"

	"github.com/LzKelvinSoares/lzk-go-lib/environment"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func Start() {
	serverVariables := environment.GetServerVariables()

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
	log.Fatal(http.ListenAndServe(serverVariables.GetAddressToServe(), router))

}

func getDBClient() *sqlx.DB {
	dbVariables := environment.GetDBVariables()

	client, err := sqlx.Open("mysql", dbVariables.GetDataSource())
	if err != nil {
		panic(err)
	}
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}
