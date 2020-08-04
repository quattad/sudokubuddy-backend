package crud

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/utils/channels"
)

// BoardsCRUDService is a global variable that exposes the methods of the module
var BoardsCRUDService BoardsCRUDInterface

func init() {
	BoardsCRUDService = &BoardsCRUD{}
}

// BoardsCRUD is a struct that makes it easy to access the db for that particular repository
// by calling r.db
type BoardsCRUD struct {
	db *gorm.DB
}

// BoardsCRUDInterface is an interface for BoardsCRUD struct to allow for mocking of functions
// during testing
type BoardsCRUDInterface interface {
	// InitDB
	NewBoardsCRUD(*gorm.DB) *BoardsCRUD

	// Create
	Save(models.Board) (models.Board, error)

	// Read
	FindByID(uint32) (models.Board, error)
	FindByPuzzleIDRowCol(uint32, int, int) (models.Board, error)
	FindAll(uint32) ([]models.Board, error)

	// Update
	Update(uint32, models.Board) (int64, error)

	// Delete
	Delete(uint32) (int64, error)
}

// NewBoardsCRUD takes in db as an argument and returns a RepositoryBoardsCRUD struct that
// has r.db as a property; making it easy to access the db
func (boardsCRUD *BoardsCRUD) NewBoardsCRUD(db *gorm.DB) *BoardsCRUD {
	boardsCRUD.db = db
	return boardsCRUD
}

// ========== CREATE ========== //

// Save takes a Board model and saves it to the db
// Returns the saved model and error if successful, returns empty Board instance and error if unsuccessful
// Also creates Board
func (boardsCRUD *BoardsCRUD) Save(board models.Board) (models.Board, error) {

	var err error
	done := make(chan bool)

	go func(ch chan<- bool) {
		err = boardsCRUD.db.Debug().Model(&models.Board{}).Create(&board).Error

		if err != nil {
			ch <- false
		}

		ch <- true
		close(ch)

	}(done)

	if channels.OK(done) {
		return board, nil
	}

	return models.Board{}, err

}

// ========== READ ========== //

// FindByID takes a board and fetches the model instance from the db
// Returns the saved model and error if successful, returns empty Board instance and error if unsuccessful
func (boardsCRUD *BoardsCRUD) FindByID(boardID uint32) (models.Board, error) {
	var err error
	board := models.Board{}
	done := make(chan bool)

	// If not found or any kind of error, will return false
	go func(ch chan<- bool) {
		defer close(ch)
		err = boardsCRUD.db.Debug().Model(&models.Board{}).Where("id=?", boardID).Take(&board).Error

		if err != nil {
			ch <- false
		}

		ch <- true
	}(done)

	// Board found
	if channels.OK(done) {
		return board, nil
	}

	// Board not found
	if gorm.IsRecordNotFoundError(err) {
		return board, errors.New("Board not found")
	}

	// Other errors
	return board, err

}

// FindByPuzzleIDRowCol takes a board and fetches the model instance from the db
// Returns the saved model and error if successful, returns empty Board instance and error if unsuccessful
func (boardsCRUD *BoardsCRUD) FindByPuzzleIDRowCol(puzzleID uint32, boardRow int, boardCol int) (models.Board, error) {
	var err error
	board := models.Board{}
	done := make(chan bool)

	// If not found or any kind of error, will return false
	go func(ch chan<- bool) {
		defer close(ch)
		err = boardsCRUD.db.Debug().Model(&models.Board{}).Where("puzzle_id=? AND board_row=? AND board_col=?", puzzleID, boardRow, boardCol).Take(&board).Error

		if err != nil {
			ch <- false
		}

		ch <- true
	}(done)

	// Board found
	if channels.OK(done) {
		return board, nil
	}

	// Board not found
	if gorm.IsRecordNotFoundError(err) {
		return board, errors.New("Board not found")
	}

	// Other errors
	return board, err

}

// FindAll fetches all the entries from the Board table by
// a. SELECT ID FROM puzzles WHERE user_id=userID (as subquery)
// b. SELECT * FROM boards WHERE puzzle_id = (result from subquery a)
// Returns an array of models and error if successful, returns empty array and error if unsuccessful
func (boardsCRUD *BoardsCRUD) FindAll(userID uint32) ([]models.Board, error) {
	var err error
	boards := []models.Board{}
	done := make(chan bool)

	// If not found or any kind of error, will return false
	go func(ch chan<- bool) {
		defer close(ch)
		// Custom SQL
		err = boardsCRUD.db.Raw("SELECT * FROM boards WHERE puzzle_id IN (SELECT ID FROM puzzles WHERE user_id=?)", userID).Scan(&boards).Error

		// err = boardsCRUD.db.Debug().Model(&models.Board{}).Limit(100).Where("puzzle_id=?", boardsCRUD.db.Debug().Model(&models.Puzzle{}).Where("user_id=?", userID).SubQuery()).Find(&boards).Error
		// err = boardsCRUD.db.Debug().Model(&models.Board{}).Limit(100).Where("puzzle_id=?", puzzleID).Find(&boards).Error

		if err != nil {
			ch <- false
		}

		ch <- true
	}(done)

	// Board found
	if channels.OK(done) {
		return boards, nil
	}

	// Other errors
	return nil, err
}

// ========== UPDATE ========== //

// Update takes in a model instance with updated fields and ID, and updates the existing entry in the db that matches ID
// Returns updated model and error if successful, returns empty array and error if unsuccessful
func (boardsCRUD *BoardsCRUD) Update(puzzleID uint32, board models.Board) (int64, error) {
	var rs *gorm.DB
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)

		rs = boardsCRUD.db.Debug().Model(&models.Board{}).Where("puzzle_id=? AND board_row=? AND board_col=?", puzzleID, board.BoardRow, board.BoardCol).UpdateColumns(
			map[string]interface{}{
				"value":      board.Value,
				"updated_at": time.Now(),
			},
		)

		ch <- true
	}(done)

	if !channels.OK(done) || rs.Error != nil {
		return 0, rs.Error
	}

	return rs.RowsAffected, nil

}

// ========== DELETE ========== //

// Delete takes in an ID and deletes the existing entry in the db that matches ID
// Returns updated model and error if successful, returns empty array and error if unsuccessful
func (boardsCRUD *BoardsCRUD) Delete(boardID uint32) (int64, error) {
	var rs *gorm.DB
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)

		rs = boardsCRUD.db.Debug().Model(&models.Board{}).Where("id=?", boardID).Delete(&models.Board{})
		ch <- true

	}(done)

	if channels.OK(done) && rs.Error == nil {
		return rs.RowsAffected, nil
	}

	return 0, rs.Error

}
