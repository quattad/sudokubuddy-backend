package database

import (
	// database/sql package must be used together with db driver. import only for side effects
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// DBService exposes the methods of dbService and allows for mocking
var DBService dbServiceInterface

func init() {
	DBService = &dbService{}
}

type dbServiceInterface interface {
	Connect(string, string) (*gorm.DB, error)
}

type dbService struct{}

// Connect establishes a connection a database at DBDRIVER and DBURL
func (d *dbService) Connect(DBDRIVER, DBURL string) (*gorm.DB, error) {
	db, err := gorm.Open(DBDRIVER, DBURL)

	if err != nil {
		return nil, err
	}

	return db, nil
}
