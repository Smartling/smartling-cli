package clierror

type UIError struct {
	Err       error
	Operation string
	Fields    map[string]string
}

func (e UIError) Error() string {
	return e.Err.Error()
}
