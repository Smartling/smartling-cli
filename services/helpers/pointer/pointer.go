package pointer

// PNew returns value for given pointer
func PNew[T any](v *T) T {
	if v == nil {
		return *new(T)
	}
	return *v
}
