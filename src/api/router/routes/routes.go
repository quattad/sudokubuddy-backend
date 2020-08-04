package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/quattad/sudokubuddy-backend/src/api/middlewares"
)

// Route is a struct that has the following fields
// URI
// Method
// Handler
// AuthRequired
type Route struct {
	URI          string
	Method       string
	Handler      func(http.ResponseWriter, *http.Request)
	AuthRequired bool
}

// Load appends each Route struct to an array of Routes and returns the array
// Iterate over and register to Router instance declared in router.go
func Load() []Route {
	routes := UserRoutes
	routes = append(routes, PuzzleRoutes...)
	routes = append(routes, LoginRoutes...)
	routes = append(routes, BoardRoutes...)
	return routes
}

// SetupRoutes takes in a mux.Router instance, registers each route returned by Load() and return it to the mux.Router instance
func SetupRoutes(r *mux.Router) *mux.Router {
	for _, route := range Load() {
		r.HandleFunc(route.URI, route.Handler).Methods(route.Method)
	}

	return r
}

// SetupRoutesWithMiddlewares registersmiddleware functions SetMiddlewareLogger, SetMiddlewareJSON and SetMiddlewareAuthentication
func SetupRoutesWithMiddlewares(r *mux.Router) *mux.Router {
	for _, route := range Load() {
		if route.AuthRequired {
			r.HandleFunc(route.URI,
				middlewares.SetMiddlewareLogger(
					middlewares.SetMiddlewareJSON(
						middlewares.SetMiddlewareAuthentication(route.Handler))),
			).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI,
				middlewares.SetMiddlewareLogger(
					middlewares.SetMiddlewareJSON(route.Handler),
				),
			).Methods(route.Method)
		}
	}

	return r
}
