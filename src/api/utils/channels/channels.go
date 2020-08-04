package channels

// OK takes in a channel and returns true once the done channel receives true
// else returns false
func OK(done <-chan bool) bool {
	select {
	case ok := <-done:
		if ok {
			return true
		}
	}

	return false
}
