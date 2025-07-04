package pointer

// NewP returns pointer for given value
func NewP[T any](v T) *T {
	return &v
}

// PNew returns value for given pointer
func PNew[T any](v *T) T {
	if v == nil {
		return *new(T)
	}
	return *v
}
