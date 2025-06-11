package clierror

import "fmt"

// ProjectNotFoundError is project not found error
type ProjectNotFoundError struct{}

// Error returns string representation.
func (error ProjectNotFoundError) Error() string {
	return NewError(
		fmt.Errorf(`specified project is not found`),
		`Check that speciied project is correct in --project option `+
			`and in config file as well.`,
	).Error()
}
