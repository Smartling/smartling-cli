package clierror

// UIError defines error
type UIError struct {
	Err         error
	Operation   string
	Description string
	Fields      map[string]string
}

func (e UIError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Description
}
