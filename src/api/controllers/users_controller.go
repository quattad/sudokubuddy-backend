package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/quattad/sudokubuddy-backend/src/api/auth"
	"github.com/quattad/sudokubuddy-backend/src/api/config"
	"github.com/quattad/sudokubuddy-backend/src/api/crud"
	"github.com/quattad/sudokubuddy-backend/src/api/database"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/responses"
)

// GetUser fetches a user by id
func GetUser(w http.ResponseWriter, r *http.Request) {
	/*
		1. Extract userID from route variables using mux.Vars() and convert to uint32, return status code 400 if err
		2. Connect to the DB, return status code 500 if err
		3. Create a new pointer to *UsersCRUD, return status code 500 if err
		4. Execute findByID, return status code 400 if err. Return status 200 and retrieved user if successful
	*/

	// Extract UserID from route variables
	routeVariables := mux.Vars(r)
	uid, err := strconv.ParseUint(routeVariables["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	// Connect to db
	db, err := database.DBService.Connect(config.DBDRIVER, config.DBURL)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	defer db.Close()

	repo := crud.UsersCRUDService.NewUsersCRUD(db)

	user, err := repo.FindByID(uint32(uid))

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	responses.JSON(w, http.StatusOK, user)

}

// GetUsers fetches all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	/*
		1. Connect to the DB, return status code 500 if err
		2. Create a new pointer to *UsersCRUD, return status code 500 if err
		3. Execute FindAl(), return status code 422 if err. Return status 200 and retrieved []models.User if successful.
	*/

	// Connect to db
	db, err := database.DBService.Connect(config.DBDRIVER, config.DBURL)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	defer db.Close()

	repo := crud.UsersCRUDService.NewUsersCRUD(db)

	users, err := repo.FindAll()

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	responses.JSON(w, http.StatusOK, users)
}

// CreateUser creates a user in the User resource
func CreateUser(w http.ResponseWriter, r *http.Request) {
	/*
		1. Read from request body into bytes. If err, return status code 400.
		2. Unmarshal from bytes to model. If err, return status code 400.
		3. Connect to db. If err, return status code 500.
	*/
	user := models.User{}
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	// Unmarshal body
	err = json.Unmarshal(body, &user)

	// Validate user
	user.PrepareUser()
	err = user.ValidateUser("")

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	// Connect to DB
	db, err := database.DBService.Connect(config.DBDRIVER, config.DBURL)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	defer db.Close()

	// Create new repository
	repo := crud.UsersCRUDService.NewUsersCRUD(db)

	user, err = repo.Save(user)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, user.ID))
	responses.JSON(w, http.StatusCreated, user)
}

// UpdateUser updates a user by id
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	/*
		1. Read from request body into bytes. If err, return status code 400.
		2. Unmarshal from bytes to model. If err, return status code 400.
		3. Connect to db. If err, return status code 500.
		4. Extract the tokenID and check if it matches the userID. If it does not match, return status code 201 unauthorized.
		5. Execute update. If successful, return status code 200 and number of rows updated.
	*/
	routeVariables := mux.Vars(r)
	uid, err := strconv.ParseUint(routeVariables["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	body, err := ioutil.ReadAll(r.Body) // read from request body

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	db, err := database.DBService.Connect(config.DBDRIVER, config.DBURL)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	defer db.Close()

	tokenUID, err := auth.TokenService.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
	}

	if tokenUID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, err)
	}

	repo := crud.UsersCRUDService.NewUsersCRUD(db)

	rows, err := repo.Update(uint32(uid), user)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	responses.JSON(w, http.StatusOK, rows)

}

// DeleteUser deletes a user by id
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	/*
		1. Extract UID from route variable
		2. Connect to db. If err, return status code 500.
		3. Extract the tokenID and check if it matches the userID. If it does not match, return status code 201 unauthorized.
		4. Execute update. If successful, return status code 200 and number of rows updated.
	*/

	routeVariables := mux.Vars(r)
	uid, err := strconv.ParseUint(routeVariables["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	db, err := database.DBService.Connect(config.DBDRIVER, config.DBURL)
	defer db.Close()

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	tokenUID, err := auth.TokenService.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
	}

	if tokenUID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, err)
	}

	repo := crud.UsersCRUDService.NewUsersCRUD(db)

	rows, err := repo.Delete(uint32(uid))

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	responses.JSON(w, http.StatusOK, rows)

}
