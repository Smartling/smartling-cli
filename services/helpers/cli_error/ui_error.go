package clierror

// UIError defines error
type UIError struct {
	Err         error
	Operation   string
	Description string
	Fields      map[string]string
}

func (e UIError) Error() string {
	return e.Err.Error()
}
