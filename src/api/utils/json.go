package utils

import (
	"encoding/json"
	"io"
	"reflect"
)

// JSONEqual compares JSON from two Readers
func JSONEqual(a, b io.Reader) bool {
	var first, second interface{}

	// Encoding, decoding is applying structure to stream of bytes
	// Marshaling, demarshaling is applying structure to in-mem bytes
	d := json.NewDecoder(a)

	if err := d.Decode(&first); err != nil {
		return false
	}

	d = json.NewDecoder(b)

	if err := d.Decode(&second); err != nil {
		return false
	}

	return reflect.DeepEqual(first, second)
}
