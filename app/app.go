package app

import (
	"banking-auth/constants"
	"banking-auth/domain"
	"banking-auth/service"
	"log"
	"net/http"

	libConstants "github.com/LzKelvinSoares/lzk-go-lib/constants"
	"github.com/LzKelvinSoares/lzk-go-lib/environment"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func Start() {
	serverVariables := environment.GetServerVariables()

	router := mux.NewRouter()
	dbType := libConstants.SQLServer
	dbClient := environment.GetSQLClient(dbType)

	//define handlers
	ah := AuthHandler{
		service: service.NewDefaultAuthService(GetDBRepositories(dbType, dbClient), domain.GetRolePermissions()),
	}

	// define routes
	router.HandleFunc(constants.LOGIN, ah.Login).Methods(http.MethodPost)
	router.HandleFunc(constants.VERIFY, ah.Verify).Methods(http.MethodGet)
	router.HandleFunc(constants.REFRESH, ah.Refresh).Methods(http.MethodPost)

	// starting server
	log.Fatal(http.ListenAndServe(serverVariables.GetAddressToServe(), router))

}

func GetDBRepositories(dbType libConstants.DBType, dbClient *sqlx.DB) domain.AuthRepository {
	var authRepository domain.AuthRepository

	switch dbType {
	case libConstants.MySQL:
		authRepository = domain.NewAuthRepositoryDB(dbClient)
	case libConstants.SQLServer:
		authRepository = domain.NewAuthRepositorySQLServerDB(dbClient)
	}

	return authRepository
}
