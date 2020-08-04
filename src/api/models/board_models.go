package models

import (
	"errors"
	"strings"
	"time"
)

// Board is a struct that defines fields in the db
type Board struct {
	ID        uint32    `gorm:"primary_key;auto_increment;unique" json:"id"`
	BoardRow  int       `gorm:"type:tinyint(1) unsigned; not null" json:"board_row"`
	BoardCol  int       `gorm:"type:tinyint(1) unsigned; not null" json:"board_col"`
	Value     int       `gorm:"type:tinyint(1) unsigned; default:0; not null" json:"value"`
	CreatedAt time.Time `gorm:"default:current_timestamp()" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp()" json:"updated_at"`
	PuzzleID  uint32    `gorm:"not null" json:"puzzle_id"`
}

// PrepareBoard removes whitespaces from puzzle fields and populates
// CreatedAt, UpdatedAt columns
func (board *Board) PrepareBoard() *Board {
	board.CreatedAt = time.Now()
	board.UpdatedAt = time.Now()
	return board
}

// ValidateBoard checks if any of the above fields are empty
func (board *Board) ValidateBoard(action string) error {
	var err error

	switch strings.ToLower(action) {
	case "update":
		if board.BoardRow < 1 {
			return errors.New("Board has invalid value for property 'board_row'")
		}
		if board.BoardCol < 1 {
			return errors.New("Board has invalid value for property 'board_col'")
		}
		if board.Value < 1 {
			return errors.New("Board has invalid value for property 'value'")
		}
		if board.PuzzleID < 1 {
			return errors.New("Board has invalid value for property 'puzzle_id'")
		}
	default:
		if board.BoardRow < 1 {
			return errors.New("Board has invalid value for property 'board_row'")
		}
		if board.BoardCol < 1 {
			return errors.New("Board has invalid value for property 'board_col'")
		}
		if board.PuzzleID < 1 {
			return errors.New("Board has invalid value for property 'puzzle_id'")
		}
	}

	err = nil
	return err
}
