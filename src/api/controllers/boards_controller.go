package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/quattad/sudokubuddy-backend/src/api/auth"
	"github.com/quattad/sudokubuddy-backend/src/api/config"
	"github.com/quattad/sudokubuddy-backend/src/api/crud"
	"github.com/quattad/sudokubuddy-backend/src/api/database"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/responses"
)

// GetBoard fetches a board by boardID
func GetBoard(w http.ResponseWriter, r *http.Request) {
	/*
		1. Extract id (boardID) from route variables using mux.Vars() and convert to uint32, return status code 400 if err
		2. Get uid (userID) from request, if not authorized, return status code 201
		3. Connect to the DB, return status code 500 if err
		4. Create a new pointer to *BoardsCRUD, return status code 500 if err
		5. Execute FindByID, return status code 400 if err.
		6. Fetch board, get corresponding puzzle with board.PuzzleID, check puzzle.UserID == uid, if not match return status code 201
		6. Return status 200 and retrieved board if successful
	*/

	// Extract ID from route variables
	routeVariables := mux.Vars(r)
	boardID, err := strconv.ParseUint(routeVariables["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	// Fetch user ID from request body
	uid, err := auth.TokenService.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
	}

	// Connect to db
	db, err := database.DBService.Connect(config.DBDRIVER, config.DBURL)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	defer db.Close()

	repo := crud.BoardsCRUDService.NewBoardsCRUD(db)

	board, err := repo.FindByID(uint32(boardID))

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	// Search board.PuzzleID to retrieve corresponding puzzle and check
	// retrieved puzzle.UserID with uid extracted from request body
	// If no match, return 201
	repoPuzzles := crud.PuzzlesCRUDService.NewPuzzlesCRUD(db)

	puzzle, err := repoPuzzles.FindByID(board.PuzzleID, uid)

	if puzzle.UserID != uid {
		responses.ERROR(w, http.StatusUnauthorized, err)
	}

	responses.JSON(w, http.StatusOK, board)
}

// GetBoardsByPuzzleIDRowCol fetches boards by puzzleID, row and column
func GetBoardsByPuzzleIDRowCol(w http.ResponseWriter, r *http.Request) {
	/*
		1. Extract puzzleID, row and col from route variables using mux.Vars() and convert to uint32, return status code 400 if err
		2. Get uid (userID) from request, if not authorized, return status code 201
		3. Connect to the DB, return status code 500 if err
		4. Create a new pointer to *BoardsCRUD, return status code 500 if err
		5. Execute FindByID, return status code 400 if err.
		6. Fetch board, get corresponding puzzle with board.PuzzleID, check puzzle.UserID == uid, if not match return status code 201
		6. Return status 200 and retrieved board if successful
	*/

	// Extract ID from route variables
	routeVariables := mux.Vars(r)
	puzzleID, err := strconv.ParseUint(routeVariables["puzzle_id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	boardRow, err := strconv.ParseUint(routeVariables["board_row"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	boardCol, err := strconv.ParseUint(routeVariables["board_col"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	// Fetch user ID from request body
	uid, err := auth.TokenService.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
	}

	// Connect to db
	db, err := database.DBService.Connect(config.DBDRIVER, config.DBURL)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	defer db.Close()

	repo := crud.BoardsCRUDService.NewBoardsCRUD(db)

	board, err := repo.FindByPuzzleIDRowCol(uint32(puzzleID), int(boardRow), int(boardCol))

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	// Search board.PuzzleID to retrieve corresponding puzzle and check
	// retrieved puzzle.UserID with uid extracted from request body
	// If no match, return 201
	repoPuzzles := crud.PuzzlesCRUDService.NewPuzzlesCRUD(db)

	puzzle, err := repoPuzzles.FindByID(board.PuzzleID, uid)

	if puzzle.UserID != uid {
		responses.ERROR(w, http.StatusUnauthorized, err)
	}

	responses.JSON(w, http.StatusOK, board)
}

// GetBoards fetches all boards
func GetBoards(w http.ResponseWriter, r *http.Request) {
	/*
		1. Get uid (userID) from request, if not authorized, return status code 201
		2. Connect to the DB, return status code 500 if err
		3. Create a new pointer to *boardsCRUD, return status code 500 if err
		4. Execute FindAll, return status code 400 if err.
		5. Check if puzzles.UserID == uid, if not match return status code 201
		6. Return status 200 and retrieved puzzles if successful
	*/

	// Fetch user ID from request body
	uid, err := auth.TokenService.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
	}

	// Connect to db
	db, err := database.DBService.Connect(config.DBDRIVER, config.DBURL)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	defer db.Close()

	repo := crud.BoardsCRUDService.NewBoardsCRUD(db)
	boards, err := repo.FindAll(uid)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	responses.JSON(w, http.StatusOK, boards)
}

// CreateBoard creates a board in the Board resource
func CreateBoard(w http.ResponseWriter, r *http.Request) {
	/*
		1. Read from request body into bytes. If err, return status code 400.
		2. Unmarshal from bytes to model. If err, return status code 400.
		3. Connect to db. If err, return status code 500.
	*/
	board := models.Board{}
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	// Unmarshal body
	err = json.Unmarshal(body, &board)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	// uid, err := auth.TokenService.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
	}

	// TODO - Resolve checking of userID
	// board.UserID = uid

	// if uid != board.UserID {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("UserID does not match the board UserID"))
	// 	return
	// }

	// Validate board
	board.PrepareBoard()
	err = board.ValidateBoard("") // default case

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
	repo := crud.BoardsCRUDService.NewBoardsCRUD(db)

	board, err = repo.Save(board)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, board.ID))
	responses.JSON(w, http.StatusCreated, board)
}

// UpdateBoard updates a board by id
func UpdateBoard(w http.ResponseWriter, r *http.Request) {
	/*
		1. Read from request body into bytes. If err, return status code 400.
		2. Unmarshal from bytes to model. If err, return status code 400.
		3. Extract tokenID from request and check if it matches userID. If no match, return status code 401.
		4. Connect to db. If err, return status code 500.
		5. Execute update. If successful, return status code 200 with number of rows updated.
	*/

	// Extract id (boardID) from route variables
	// For testing purposes
	// routeVariables := mux.Vars(r)
	// puzzleID, err := strconv.ParseUint(routeVariables["puzzle_id"], 10, 32)

	// if err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// }

	// boardRow, err := strconv.ParseUint(routeVariables["board_row"], 10, 32)

	// if err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// }

	// boardCol, err := strconv.ParseUint(routeVariables["board_col"], 10, 32)

	// if err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// }

	// Test parsing URL
	u, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	// Read from body and unmarshal into empty models.Board
	body, err := ioutil.ReadAll(r.Body) // read from request body

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	board := models.Board{}
	err = json.Unmarshal(body, &board)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	if len(u["puzzle_id"]) > 1 || len(u["board_row"]) > 1 || len(u["board_col"]) > 1 {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	// Assign BoardRow, BoardCol and Value to Board instance
	// Convert both from uint64
	puzzleID, err := strconv.Atoi(u["puzzle_id"][0])

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	board.BoardRow, err = strconv.Atoi(u["board_row"][0])

	if err != nil || board.BoardRow < 1 || board.BoardRow > 9 {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	board.BoardCol, err = strconv.Atoi(u["board_col"][0])

	if err != nil || board.BoardCol < 1 || board.BoardCol > 9 {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	// TODO - Resolve checking of userID
	// Extract tokenID from request and check if it matches userID. If no match, return status code 401.
	// userID, err := auth.TokenService.ExtractTokenID(r)

	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, err)
	// }

	// Connect to database
	db, err := database.DBService.Connect(config.DBDRIVER, config.DBURL)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	defer db.Close()

	// Execute search
	repo := crud.BoardsCRUDService.NewBoardsCRUD(db)

	rows, err := repo.Update(uint32(puzzleID), board)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	responses.JSON(w, http.StatusOK, rows)
}

// DeleteBoard deletes a board by id
func DeleteBoard(w http.ResponseWriter, r *http.Request) {
	/*
		1. Extract UID from route variable
		2. Connect to db. If err, return status code 500.
		3. Extract the tokenID and check if it matches the userID. If it does not match, return status code 201 unauthorized.
		4. Execute update. If successful, return status code 200 and number of rows updated.
	*/

	routeVariables := mux.Vars(r)
	boardID, err := strconv.ParseUint(routeVariables["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	db, err := database.DBService.Connect(config.DBDRIVER, config.DBURL)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	defer db.Close()

	// tokenUID, err := auth.TokenService.ExtractTokenID(r)

	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, err)
	// }

	repo := crud.BoardsCRUDService.NewBoardsCRUD(db)

	rows, err := repo.Delete(uint32(boardID))

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	responses.JSON(w, http.StatusOK, rows)
}
