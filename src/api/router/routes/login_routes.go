package routes

import (
	"net/http"

	"github.com/quattad/sudokubuddy-backend/src/api/controllers"
)

// LoginRoutes is an array of Route instances which map paths to route handlers
var LoginRoutes = []Route{
	Route{
		URI:          "/login",
		Method:       http.MethodPost,
		Handler:      controllers.LoginControllerService.Login,
		AuthRequired: false,
	},
}
