package routes

import (
	"net/http"

	"github.com/quattad/sudokubuddy-backend/src/api/controllers"
)

// UserRoutes is an array of Route structs with URI, Method, Handler and AuthRequired fields
var UserRoutes = []Route{
	Route{
		URI:          "/users",
		Method:       http.MethodGet,
		Handler:      controllers.GetUsers,
		AuthRequired: false,
	},
	Route{
		URI:          "/users/{id}",
		Method:       http.MethodGet,
		Handler:      controllers.GetUser,
		AuthRequired: false,
	},
	Route{
		URI:          "/users",
		Method:       http.MethodPost,
		Handler:      controllers.CreateUser,
		AuthRequired: false,
	},
	Route{
		URI:          "/users/{id}",
		Method:       http.MethodPut,
		Handler:      controllers.UpdateUser,
		AuthRequired: true,
	},
	Route{
		URI:          "/users/{id}",
		Method:       http.MethodDelete,
		Handler:      controllers.DeleteUser,
		AuthRequired: true,
	},
}
