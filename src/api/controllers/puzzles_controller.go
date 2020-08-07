package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"github.com/quattad/sudokubuddy-backend/src/api/auth"
	"github.com/quattad/sudokubuddy-backend/src/api/caching"
	"github.com/quattad/sudokubuddy-backend/src/api/config"
	"github.com/quattad/sudokubuddy-backend/src/api/crud"
	"github.com/quattad/sudokubuddy-backend/src/api/database"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/responses"
)

// GetPuzzle fetches a puzzle by id and user_id
func GetPuzzle(w http.ResponseWriter, r *http.Request) {
	/*
		1. Extract id (puzzleID) from route variables using mux.Vars() and convert to uint32, return status code 400 if err
		2. Get uid (userID) from request, if not authorized, return status code 201
		3. Connect to the DB, return status code 500 if err
		4. Create a new pointer to *PuzzlesCRUD, return status code 500 if err
		5. Execute FindByID, return status code 400 if err.
		6. Check if puzzles.UserID == uid, if not match return status code 201
		6. Return status 200 and retrieved puzzle if successful
	*/

	// Extract ID from route variables
	routeVariables := mux.Vars(r)
	pid, err := strconv.ParseUint(routeVariables["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	// Fetch user ID from request body
	uid, err := auth.TokenService.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
	}

	if it, found := caching.Cache.Get("puzzles/" + strconv.Itoa(int(pid))); found {

		puzzle := models.Puzzle{}

		err = json.Unmarshal(it.([]byte), &puzzle)

		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
		}

		responses.JSON(w, http.StatusOK, puzzle)

	} else {

		// Connect to db
		db, err := database.DBService.Connect(config.DBDRIVER, config.DBURL)

		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
		}

		defer db.Close()

		repo := crud.PuzzlesCRUDService.NewPuzzlesCRUD(db)

		puzzle, err := repo.FindByID(uint32(pid), uid)

		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
		}

		b, err := json.Marshal(puzzle)

		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
		}

		// GOCACHE
		caching.Cache.Set("puzzles/"+strconv.Itoa(int(pid)), b, cache.DefaultExpiration)

		responses.JSON(w, http.StatusOK, puzzle)
	}
}

// GetPuzzles fetches all puzzles
func GetPuzzles(w http.ResponseWriter, r *http.Request) {
	/*
		1. Get uid (userID) from request, if not authorized, return status code 201
		2. Connect to the DB, return status code 500 if err
		3. Create a new pointer to *PuzzlesCRUD, return status code 500 if err
		4. Execute FindAll, return status code 400 if err.
		5. Check if puzzles.UserID == uid, if not match return status code 201
		6. Return status 200 and retrieved puzzles if successful
	*/

	// Fetch user ID from request body
	uid, err := auth.TokenService.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
	}

	if it, found := caching.Cache.Get("puzzles/all"); found {

		var puzzles []models.Puzzle
		err := json.Unmarshal(it.([]byte), &puzzles)

		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
		}

		responses.JSON(w, http.StatusOK, puzzles)

	} else {

		// Connect to db
		db, err := database.DBService.Connect(config.DBDRIVER, config.DBURL)

		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
		}

		defer db.Close()

		repo := crud.PuzzlesCRUDService.NewPuzzlesCRUD(db)

		puzzles, err := repo.FindAll(uid)

		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
		}

		b, err := json.Marshal(puzzles)

		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
		}

		// GOCACHE
		caching.Cache.Set("puzzles/all", b, cache.DefaultExpiration)

		responses.JSON(w, http.StatusOK, puzzles)
	}
}

// CreatePuzzle creates a puzzle in the Puzzle resource
func CreatePuzzle(w http.ResponseWriter, r *http.Request) {
	/*
		1. Read from request body into bytes. If err, return status code 400.
		2. Unmarshal from bytes to model. If err, return status code 400.
		3. Connect to db. If err, return status code 500.
	*/
	puzzle := models.Puzzle{}
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	// Unmarshal body
	err = json.Unmarshal(body, &puzzle)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	uid, err := auth.TokenService.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
	}

	puzzle.UserID = uid

	// if uid != puzzle.UserID {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("UserID does not match the puzzle UserID"))
	// 	return
	// }

	// Validate puzzle
	puzzle.PreparePuzzle()
	err = puzzle.ValidatePuzzle("") // default case

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
	repo := crud.PuzzlesCRUDService.NewPuzzlesCRUD(db)

	puzzle, err = repo.Save(puzzle)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, puzzle.ID))
	responses.JSON(w, http.StatusCreated, puzzle)
}

// UpdatePuzzle updates a puzzle by id
func UpdatePuzzle(w http.ResponseWriter, r *http.Request) {
	/*
		1. Read from request body into bytes. If err, return status code 400.
		2. Unmarshal from bytes to model. If err, return status code 400.
		3. Extract tokenID from request and check if it matches userID. If no match, return status code 401.
		4. Connect to db. If err, return status code 500.
		5. Execute update. If successful, return status code 200 with number of rows updated.
	*/

	// Extract id (puzzleID) from route variables
	routeVariables := mux.Vars(r)
	puzzleID, err := strconv.ParseUint(routeVariables["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	// Read from body and unmarshal into empty models.Puzzle
	body, err := ioutil.ReadAll(r.Body) // read from request body

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	puzzle := models.Puzzle{}
	err = json.Unmarshal(body, &puzzle)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	// Write puzzleID to puzzle model
	puzzle.ID = uint32(puzzleID)

	// Extract tokenID from request and check if it matches userID. If no match, return status code 401.
	userID, err := auth.TokenService.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
	}

	// Write userID to puzzle model
	puzzle.UserID = userID

	// Connect to database
	db, err := database.DBService.Connect(config.DBDRIVER, config.DBURL)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	defer db.Close()

	// Execute search
	repo := crud.PuzzlesCRUDService.NewPuzzlesCRUD(db)

	rows, err := repo.Update(userID, puzzle)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	responses.JSON(w, http.StatusOK, rows)
}

// DeletePuzzle deletes a puzzle by id
func DeletePuzzle(w http.ResponseWriter, r *http.Request) {
	/*
		1. Extract UID from route variable
		2. Connect to db. If err, return status code 500.
		3. Extract the tokenID and check if it matches the userID. If it does not match, return status code 201 unauthorized.
		4. Execute update. If successful, return status code 200 and number of rows updated.
	*/

	routeVariables := mux.Vars(r)
	puzzleID, err := strconv.ParseUint(routeVariables["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
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

	repo := crud.PuzzlesCRUDService.NewPuzzlesCRUD(db)

	rows, err := repo.Delete(uint32(puzzleID), tokenUID)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	responses.JSON(w, http.StatusOK, rows)
}
