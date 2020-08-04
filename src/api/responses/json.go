package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSON is a function
// httpResponseWriter - interface used by HTTP handler to create HTTP response
// httpResponseWriter.WriteHeader(statusCode int) - sends HTTP response header with provided status code
// json.NewEncoder(w) - returns new encoder that writes to w
// type Encoder - writes JSON values to output stream
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	/*
		1. Writes the status code to the response header
		2. Create a new Encoder that writes to w and encode data to w, return error if unsuccessful
	*/
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}

}

// ERROR formats a valid response object with the following format
//
func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		data := struct {
			Error string `json:"error"` // `json:"error" specifies json key as 'error'`
		}{
			Error: err.Error(),
		}

		JSON(w, statusCode, data)
	}

	JSON(w, http.StatusBadRequest, nil)
}
