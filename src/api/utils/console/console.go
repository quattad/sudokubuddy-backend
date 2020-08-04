package console

import (
	"encoding/json"
	"fmt"
	"log"
)

// Pretty formats model to be printed out
func Pretty(data interface{}) {
	d, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(d))
}
