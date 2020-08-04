package router

import (
	"github.com/gorilla/mux"
	"github.com/quattad/sudokubuddy-backend/src/api/router/routes"
)

// New creates a router instanceand sets up paths with the router instance
func New() *mux.Router {

	// Strictslash defines trailing slash behaviour
	// e.g. /posts/ will be redirected to /posts
	r := mux.NewRouter().StrictSlash(true)
	return routes.SetupRoutesWithMiddlewares(r)
	// return routes.SetupRoutes(r)
}
