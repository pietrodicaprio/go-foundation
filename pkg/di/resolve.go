package di

// Resolve retrieves a typed dependency from the container.
// Returns zero value and false if not found or type mismatch.
func Resolve[T any](c *Container, name string) (T, bool) {
	var zero T
	v, ok := c.Get(name)
	if !ok {
		return zero, false
	}
	typed, ok := v.(T)
	if !ok {
		return zero, false
	}
	return typed, true
}

// MustResolve retrieves a typed dependency and panics if not found.
func MustResolve[T any](c *Container, name string) T {
	v, ok := Resolve[T](c, name)
	if !ok {
		panic("di: cannot resolve " + name)
	}
	return v
}
