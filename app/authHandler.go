package app

import (
	"banking-auth/dto"
	"banking-auth/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AuthHandler struct {
	service service.AuthService
}

func (ah AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		log.Println("Error while decoding login request: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, logErr := ah.service.Login(loginRequest)
		if logErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, logErr.Error())
		} else {
			fmt.Fprint(w, *token)
		}
	}
}
