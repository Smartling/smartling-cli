package clierror

import "fmt"

type Error struct {
	Cause       error
	Description string
	Args        []interface{}
}

func NewError(cause error, description string, args ...interface{}) Error {
	return Error{
		Cause:       cause,
		Description: description,
		Args:        args,
	}
}

func (err Error) Unwrap() error {
	return err.Cause
}

func (err Error) Error() string {
	return fmt.Sprintf(
		"ERROR: %s\n\n%s",
		err.Cause,
		fmt.Sprintf(err.Description, err.Args...),
	)
}

type ProjectNotFoundError struct{}

func (error ProjectNotFoundError) Error() string {
	return NewError(
		fmt.Errorf(`specified project is not found`),
		`Check that speciied project is correct in --project option `+
			`and in config file as well.`,
	).Error()
}
