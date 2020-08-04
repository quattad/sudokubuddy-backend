package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/quattad/sudokubuddy-backend/src/api/auth"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/responses"
)

// LoginControllerService is a global variable that exposes the login controller module
var (
	LoginControllerService loginControllerInterface
)

func init() {
	LoginControllerService = &loginControllerService{}
}

type loginControllerService struct{}

type loginControllerInterface interface {
	Login(http.ResponseWriter, *http.Request)
}

// Login authenticates a user. Returns a token if login is successful, returns an error if login is unsuccessful
func (l *loginControllerService) Login(w http.ResponseWriter, r *http.Request) {
	// ReadAll reads from r until error or EOF is encountered, then returns
	// r is io.Reader type
	// https://golang.org/pkg/io/ioutil/#ReadAll
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Write details of
	user := models.User{}
	err = json.Unmarshal(body, &user)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.PrepareUser()               // remove whitespaces
	err = user.ValidateUser("login") // sanitizes the fields

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// generates token if login is successful
	// err != nil if login is unsuccessful
	token, err := auth.AuthService.SignIn(user.Email, user.Password)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	responses.JSON(w, http.StatusOK, token)

}
