package main

import (
	"encoding/json"
	"github.com/reconquest/hierr-go"
)

type FormatExecutionError struct {
	Cause  error
	Format string
	Data   interface{}
}

func (err FormatExecutionError) Error() string {
	data, _ := json.MarshalIndent(err.Data, "", "  ")

	return NewError(
		hierr.Push(
			"template execution failed",
			hierr.Push(
				"error",
				err.Cause,
			),
			hierr.Push(
				"template",
				err.Format,
			),
			hierr.Push(
				"data given to template",
				data,
			),
		),

		"Data that was given to the template can't match template "+
			"definition.\n\nCheck that all fields in given data match "+
			"template.",
	).Error()
}
