package routes

import (
	"net/http"

	"github.com/quattad/sudokubuddy-backend/src/api/controllers"
)

// PuzzleRoutes is an array of Route instances which map paths to route handlers
var PuzzleRoutes = []Route{
	Route{
		URI:          "/puzzles",
		Method:       http.MethodGet,
		Handler:      controllers.GetPuzzles,
		AuthRequired: true,
	},
	Route{
		URI:          "/puzzles/{id}",
		Method:       http.MethodGet,
		Handler:      controllers.GetPuzzle,
		AuthRequired: true,
	},
	Route{
		URI:          "/puzzles",
		Method:       http.MethodPost,
		Handler:      controllers.CreatePuzzle,
		AuthRequired: true,
	},
	Route{
		URI:          "/puzzles/{id}",
		Method:       http.MethodPut,
		Handler:      controllers.UpdatePuzzle,
		AuthRequired: true,
	},
	Route{
		URI:          "/puzzles/{id}",
		Method:       http.MethodDelete,
		Handler:      controllers.DeletePuzzle,
		AuthRequired: true,
	},
}
