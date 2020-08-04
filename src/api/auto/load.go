package auto

import (
	"log"

	"github.com/quattad/sudokubuddy-backend/src/api/config"
	"github.com/quattad/sudokubuddy-backend/src/api/database"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
)

// Load attempts connection to the database
func Load() {
	db, err := database.DBService.Connect(config.DBDRIVER, config.DBURL)

	if err != nil {
		// print error, followed by call to os.exit
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Debug().DropTableIfExists(&models.Board{}, &models.Puzzle{}, &models.User{}).Error

	if err != nil {
		// print error, followed by call to os.exit
		log.Fatal(err)
	}

	// Creates tables for models based on schema defined in models
	err = db.Debug().AutoMigrate(&models.Board{}, &models.Puzzle{}, &models.User{}).Error

	if err != nil {
		// print error, followed by call to os.exit
		log.Fatal(err)
	}

	// Add foreign key constraints
	// ID in User model(PK) - UserID in Puzzle model (FK)
	err = db.Debug().Model(&models.Puzzle{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error

	if err != nil {
		// print error, followed by call to os.exit
		log.Fatal(err)
	}

	// ID in Puzzle model(PK) - PuzzleID in Board model (FK)
	err = db.Debug().Model(&models.Board{}).AddForeignKey("puzzle_id", "puzzles(id)", "cascade", "cascade").Error

	if err != nil {
		log.Fatal(err)
	}

	// Populate db with initial values
	for i, _ := range users {

		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error

		if err != nil {
			log.Fatal(err)
		}
	}

	for j, _ := range puzzles {
		// puzzles[j].UserID = users[i].ID

		err = db.Debug().Model(&models.Puzzle{}).Create(&puzzles[j]).Error

		if err != nil {
			log.Fatal(err)
		}

		// Create board for every new puzzle
		for k := 1; k <= 9; k++ {
			for l := 1; l <= 9; l++ {

				board := models.Board{}
				board.BoardRow = k
				board.BoardCol = l
				board.PuzzleID = puzzles[j].ID

				err = db.Debug().Model(&models.Board{}).Create(&board).Error

				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	// 	console.Pretty(posts[i])
	// }
}
