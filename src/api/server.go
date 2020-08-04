package api

import (
	"fmt"
	"net/http"

	"github.com/quattad/sudokubuddy-backend/src/api/auto"
	"github.com/quattad/sudokubuddy-backend/src/api/config"
	"github.com/quattad/sudokubuddy-backend/src/api/router"
)

// Run runs a server instance
func Run() {
	config.Load()
	auto.Load()
	fmt.Printf("\n\t Listening on PORT:%d\n", config.PORT) // to replace PORT with config
	Listen(config.PORT)
}

// Listen initializes a new Router instance using the mux package
func Listen(port int) {
	r := router.New()
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
