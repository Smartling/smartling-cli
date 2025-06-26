package clierror

import "fmt"

// Error is a custom error type for CLI applications.
type Error struct {
	Cause       error
	Description string
	Args        []interface{}
}

// NewError creates a new Error instance with the provided cause, description, and optional arguments.
func NewError(cause error, description string, args ...interface{}) Error {
	return Error{
		Cause:       cause,
		Description: description,
		Args:        args,
	}
}

// Unwrap returns cause of the error.
func (err Error) Unwrap() error {
	return err.Cause
}

// Error returns string representation.
func (err Error) Error() string {
	return fmt.Sprintf(
		"ERROR: %s\n\n%s",
		err.Cause,
		fmt.Sprintf(err.Description, err.Args...),
	)
}
