package utils

// Must panics on non-nil errors.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
