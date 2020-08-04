package crud

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/utils/channels"
)

// PuzzlesCRUDService is a global variable that exposes the methods of the module
var PuzzlesCRUDService PuzzlesCRUDInterface

func init() {
	PuzzlesCRUDService = &PuzzlesCRUD{}
}

// PuzzlesCRUD is a struct that makes it easy to access the db for that particular repository
// by calling r.db
type PuzzlesCRUD struct {
	db *gorm.DB
}

// PuzzlesCRUDInterface is an interface for PuzzlesCRUD struct to allow for mocking of functions
// during testing
type PuzzlesCRUDInterface interface {
	// InitDB
	NewPuzzlesCRUD(*gorm.DB) *PuzzlesCRUD

	// Create
	Save(models.Puzzle) (models.Puzzle, error)

	// Read
	FindByID(uint32, uint32) (models.Puzzle, error)
	FindAll(uint32) ([]models.Puzzle, error)

	// Update
	Update(uint32, models.Puzzle) (int64, error)

	// Delete
	Delete(uint32, uint32) (int64, error)
}

// NewPuzzlesCRUD takes in db as an argument and returns a RepositoryPuzzlesCRUD struct that
// has r.db as a property; making it easy to access the db
func (puzzlesCRUD *PuzzlesCRUD) NewPuzzlesCRUD(db *gorm.DB) *PuzzlesCRUD {
	puzzlesCRUD.db = db
	return puzzlesCRUD
}

// ========== CREATE ========== //

// Save takes a Puzzle model and saves it to the db
// Returns the saved model and error if successful, returns empty Puzzle instance and error if unsuccessful
// Also creates Board
func (puzzlesCRUD *PuzzlesCRUD) Save(puzzle models.Puzzle) (models.Puzzle, error) {

	var err error
	done := make(chan bool)

	go func(ch chan<- bool) {
		err = puzzlesCRUD.db.Debug().Model(&models.Puzzle{}).Create(&puzzle).Error

		if err != nil {
			ch <- false
		}

		// Create board for every new puzzle
		for i := 1; i <= 9; i++ {
			for j := 1; j <= 9; j++ {

				board := models.Board{}
				board.BoardRow = i
				board.BoardCol = j
				board.PuzzleID = puzzle.ID
				err = puzzlesCRUD.db.Debug().Model(&models.Board{}).Create(&board).Error

				if err != nil {
					ch <- false
				}
			}
		}

		ch <- true
		close(ch)

	}(done)

	if channels.OK(done) {
		return puzzle, nil
	}

	return models.Puzzle{}, err

}

// ========== READ ========== //

// FindByID takes a userID and fetches the model instance from the db
// Returns the saved model and error if successful, returns empty Puzzle instance and error if unsuccessful
func (puzzlesCRUD *PuzzlesCRUD) FindByID(puzzleID uint32, userID uint32) (models.Puzzle, error) {
	var err error
	puzzle := models.Puzzle{}
	done := make(chan bool)

	// If not found or any kind of error, will return false
	go func(ch chan<- bool) {
		defer close(ch)
		err = puzzlesCRUD.db.Debug().Model(&models.Puzzle{}).Where("id=? AND user_id=?", puzzleID, userID).Take(&puzzle).Error

		if err != nil {
			ch <- false
		}

		ch <- true
	}(done)

	// Puzzle found
	if channels.OK(done) {
		return puzzle, nil
	}

	// Puzzle not found
	if gorm.IsRecordNotFoundError(err) {
		return puzzle, errors.New("Puzzle not found")
	}

	// Other errors
	return puzzle, err

}

// FindAll fetches all the entries from the Puzzle model in the db
// Returns an array of models and error if successful, returns empty array and error if unsuccessful
func (puzzlesCRUD *PuzzlesCRUD) FindAll(userID uint32) ([]models.Puzzle, error) {
	var err error
	puzzles := []models.Puzzle{}
	done := make(chan bool)

	// If not found or any kind of error, will return false
	go func(ch chan<- bool) {
		defer close(ch)
		err = puzzlesCRUD.db.Debug().Model(&models.Puzzle{}).Limit(100).Where("user_id=?", userID).Find(&puzzles).Error

		if err != nil {
			ch <- false
		}

		ch <- true
	}(done)

	// Puzzle found
	if channels.OK(done) {
		return puzzles, nil
	}

	// Other errors
	return nil, err
}

// ========== UPDATE ========== //

// Update takes in a model instance with updated fields and ID, and updates the existing entry in the db that matches ID
// Returns updated model and error if successful, returns empty array and error if unsuccessful
func (puzzlesCRUD *PuzzlesCRUD) Update(userID uint32, puzzle models.Puzzle) (int64, error) {
	var rs *gorm.DB
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)

		rs = puzzlesCRUD.db.Debug().Model(&models.Puzzle{}).Where("id=? AND user_id=?", puzzle.ID, userID).UpdateColumns(
			map[string]interface{}{
				"name":       puzzle.Name,
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
func (puzzlesCRUD *PuzzlesCRUD) Delete(puzzleID uint32, userID uint32) (int64, error) {
	var rs *gorm.DB
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)

		rs = puzzlesCRUD.db.Debug().Model(&models.Puzzle{}).Where("id=? AND user_id=?", puzzleID, userID).Delete(&models.Puzzle{})
		ch <- true

	}(done)

	if channels.OK(done) && rs.Error == nil {
		return rs.RowsAffected, nil
	}

	return 0, rs.Error

}
