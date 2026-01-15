package di

import (
	"reflect"
	"sync"
)

// Container manages dependency injection with thread-safe access.
//
// Example:
//
//	c := di.New()
//	c.Provide("db", &Database{})
//	db := di.Get[*Database](c, "db")
type Container struct {
	providers map[string]any
	mu        sync.RWMutex
}

// New creates an empty DI container.
func New() *Container {
	return &Container{
		providers: make(map[string]any),
	}
}

// Provide registers a dependency by name.
func (c *Container) Provide(name string, instance any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.providers[name] = instance
}

// Get retrieves a dependency by name.
//
// Returns:
//
// The value and true if found, otherwise nil and false.
func (c *Container) Get(name string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.providers[name]
	return v, ok
}

// MustGet retrieves a dependency and panics if not found.
//
// Notes:
//
// Panics if the dependency is missing.
func (c *Container) MustGet(name string) any {
	v, ok := c.Get(name)
	if !ok {
		panic("di: dependency not found: " + name)
	}
	return v
}

// Has checks if a dependency is registered.
func (c *Container) Has(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.providers[name]
	return ok
}

// Inject populates struct fields with registered dependencies.
// Uses `inject:"name"` tags to match fields with providers.
// Falls back to field name if no inject tag is present.
func (c *Container) Inject(target any) {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return
	}

	elem := val.Elem()
	elemType := elem.Type()

	c.mu.RLock()
	defer c.mu.RUnlock()

	for i := 0; i < elem.NumField(); i++ {
		fieldType := elemType.Field(i)
		fieldVal := elem.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		name := fieldType.Tag.Get("inject")
		if name == "" {
			name = fieldType.Name
		}

		if dep, ok := c.providers[name]; ok {
			depVal := reflect.ValueOf(dep)
			if depVal.Type().AssignableTo(fieldType.Type) {
				fieldVal.Set(depVal)
			}
		}
	}
}

// Clone creates a shallow copy of the container.
func (c *Container) Clone() *Container {
	c.mu.RLock()
	defer c.mu.RUnlock()

	clone := New()
	for k, v := range c.providers {
		clone.providers[k] = v
	}
	return clone
}

// Keys returns all registered dependency names.
func (c *Container) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.providers))
	for k := range c.providers {
		keys = append(keys, k)
	}
	return keys
}
