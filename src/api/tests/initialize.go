package tests

import (
	"database/sql"
	"log"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

// Suite is a struct with the following fields
// mock - mock struct from sqlmock
// db - open test db with sqlmock
type Suite struct {
	Mock sqlmock.Sqlmock
	DB   *gorm.DB
}

// CreateSuite returns a custom Suite struct
func CreateSuite() (s Suite) {
	var (
		err error
		db  *sql.DB
	)

	db, s.Mock, err = sqlmock.New()

	if err != nil {
		log.Fatal(err)
	}

	s.DB, err = gorm.Open("mysql", db)

	if err != nil {
		log.Fatal(err)
	}

	return s
}
