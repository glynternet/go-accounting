package account

import "bytes"

// FieldError holds zero or more descriptions of things that are wrong with potential new Account items.
type FieldError []string

// Error ensures that FieldError adheres to the error interface.
func (e FieldError) Error() string {
	var errorString bytes.Buffer
	errorString.WriteString("FieldError: ")
	for i, field := range e {
		errorString.WriteString(field)
		if i < len(e)-1 {
			errorString.WriteByte(' ')
		}
	}
	return errorString.String()
}

// Equal returns true if two AccountFieldErrors contain the same error information strings in exactly the same order.
// Duplicate error information strings held within the FieldError are counted as individual error strings.
func (e FieldError) Equal(other FieldError) bool {
	if len(e) != len(other) {
		return false
	}
	for i := range e {
		if e[i] != other[i] {
			return false
		}
	}
	return true
}

// Various error strings describing possible errors with potential new Account items.
const (
	EmptyNameError = "Empty name."
)
