package tests

import "errors"

/* =================  errReader to check for ioutil.ReadAll errors ================= */
type ErrReader int

func (ErrReader) Read([]byte) (int, error) {
	return 0, errors.New("Read error")
}
