package common

import (
	"fmt"
	"testing"
)

// FatalIfError will fail a test immediately with a message if the error is non nil
func FatalIfError(t *testing.T, err error, message string) {
	if err == nil {
		return
	}
	t.Fatalf("%s: %s", message, err)
}

// FatalIfErrorf will fail a test immediately with a formatted message if the error is non nil
func FatalIfErrorf(t *testing.T, err error, format string, args ...interface{}) {
	FatalIfError(t, err, fmt.Sprintf(format, args...))
}

// ErrorIfError will fail a test immediately with a formatted message if the error is non nil
func ErrorIfError(t *testing.T, err error, message string) {
	if err == nil {
		return
	}
	t.Errorf("%s: %s", message, err)
}

// ErrorIfErrorf will fail a test immediately with a formatted message if the error is non nil
func ErrorIfErrorf(t *testing.T, err error, format string, args ...interface{}) {
	ErrorIfError(t, err, fmt.Sprintf(format, args...))
}
