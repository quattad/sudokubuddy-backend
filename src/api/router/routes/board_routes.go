package routes

import (
	"net/http"

	"github.com/quattad/sudokubuddy-backend/src/api/controllers"
)

// BoardRoutes is an array of Route instances which map paths to route handlers
var BoardRoutes = []Route{
	Route{
		URI:          "/boards",
		Method:       http.MethodGet,
		Handler:      controllers.GetBoards,
		AuthRequired: true,
	},
	Route{
		URI:          "/boards/{id}",
		Method:       http.MethodGet,
		Handler:      controllers.GetBoard,
		AuthRequired: true,
	},
	Route{
		URI:          "/boards",
		Method:       http.MethodPost,
		Handler:      controllers.CreateBoard,
		AuthRequired: true,
	},
	Route{
		URI:          "/boards",
		Method:       http.MethodPut,
		Handler:      controllers.UpdateBoard,
		AuthRequired: true,
	},
	Route{
		URI:          "/boards/{id}",
		Method:       http.MethodDelete,
		Handler:      controllers.DeleteBoard,
		AuthRequired: true,
	},
}
