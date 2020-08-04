package models

import (
	"errors"
	"html"
	"strings"
	"time"
)

// Puzzle is a struct that defines fields in the db
type Puzzle struct {
	ID        uint32    `gorm:"primary_key;auto_increment;unique" json:"id"`
	Name      string    `gorm:"size:20;not null;unique" json:"name"`
	CreatedAt time.Time `gorm:"default:current_timestamp()" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp()" json:"updated_at"`
	UserID    uint32    `gorm:"not null" json:"user_id"`
	Boards    []Board   `gorm:"foreignkey:PuzzleID association_foreignkey:ID" json:"boards"`
}

// PreparePuzzle removes whitespaces from puzzle fields
func (puzzle *Puzzle) PreparePuzzle() *Puzzle {
	puzzle.Name = html.EscapeString(strings.TrimSpace(puzzle.Name))
	puzzle.CreatedAt = time.Now()
	puzzle.UpdatedAt = time.Now()
	return puzzle
}

// ValidatePuzzle checks if any of the above fields are empty
func (puzzle *Puzzle) ValidatePuzzle(action string) error {
	var err error

	switch strings.ToLower(action) {
	case "update":
		if puzzle.Name == "" {
			return errors.New("Puzzle must have defined property 'name'")
		}

		if puzzle.UserID < 1 {
			return errors.New("Puzzle has invalid value for property 'user_id'")
		}
	default:
		if puzzle.Name == "" {
			return errors.New("Puzzle must have defined property 'name'")
		}

		if puzzle.UserID < 1 {
			return errors.New("Puzzle has invalid value for property 'user_id'")
		}
	}

	err = nil
	return err
}
