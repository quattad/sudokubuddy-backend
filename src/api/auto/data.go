package auto

import "github.com/quattad/sudokubuddy-backend/src/api/models"

// Define initial values to be populated into db based on models

var users = []models.User{
	models.User{
		Username:  "johndoe",
		Email:     "johndoe@gmail.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "123456",
	},
	models.User{
		Username:  "winstondoe",
		Email:     "winstondoe@gmail.com",
		FirstName: "Winston",
		LastName:  "Doe",
		Password:  "password123!",
	},
}

var puzzles = []models.Puzzle{
	models.Puzzle{
		ID:     1,
		Name:   "John's First Puzzle",
		UserID: 1,
	},
	models.Puzzle{
		ID:     2,
		Name:   "John's Second Puzzle",
		UserID: 1,
	},
	models.Puzzle{
		ID:     3,
		Name:   "W's First puzzle",
		UserID: 2,
	},
}
