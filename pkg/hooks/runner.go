package hooks

import (
	"context"
)

// Runner manages execution of lifecycle hooks with before/after patterns.
type Runner struct {
	discovery *Discovery
	before    map[string][]HookFunc
	after     map[string][]HookFunc
}

// HookFunc is a function called at a lifecycle event.
type HookFunc func(ctx context.Context, key string, args []any) error

// NewRunner creates a hook runner with a shared discovery instance.
func NewRunner() *Runner {
	return &Runner{
		discovery: NewDiscovery(),
		before:    make(map[string][]HookFunc),
		after:     make(map[string][]HookFunc),
	}
}

// Before registers a function to be called before a specific event.
func (r *Runner) Before(key string, fn HookFunc) {
	r.before[key] = append(r.before[key], fn)
}

// After registers a function to be called after a specific event.
func (r *Runner) After(key string, fn HookFunc) {
	r.after[key] = append(r.after[key], fn)
}

// BeforeAll registers a function to be called before any event.
func (r *Runner) BeforeAll(fn HookFunc) {
	r.before["*"] = append(r.before["*"], fn)
}

// AfterAll registers a function to be called after any event.
func (r *Runner) AfterAll(fn HookFunc) {
	r.after["*"] = append(r.after["*"], fn)
}

// Run executes before hooks, the action, and after hooks.
func (r *Runner) Run(ctx context.Context, key string, action func() error, args ...any) error {
	if err := r.runHooks(ctx, key, r.before, args); err != nil {
		return err
	}

	if err := action(); err != nil {
		return err
	}

	return r.runHooks(ctx, key, r.after, args)
}

func (r *Runner) runHooks(ctx context.Context, key string, hooks map[string][]HookFunc, args []any) error {
	// Global hooks first
	for _, fn := range hooks["*"] {
		if err := fn(ctx, key, args); err != nil {
			return err
		}
	}

	// Specific hooks
	for _, fn := range hooks[key] {
		if err := fn(ctx, key, args); err != nil {
			return err
		}
	}

	return nil
}

// Discovery returns the underlying discovery instance.
func (r *Runner) Discovery() *Discovery {
	return r.discovery
}

// Clear removes all registered hooks.
func (r *Runner) Clear() {
	r.before = make(map[string][]HookFunc)
	r.after = make(map[string][]HookFunc)
}
