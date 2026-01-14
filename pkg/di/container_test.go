package di

import (
	"testing"
)

type testDB struct {
	Name string
}

type testService struct {
	DB     *testDB `inject:"db"`
	Logger string  `inject:"logger"`
}

func TestContainer_ProvideAndGet(t *testing.T) {
	c := New()
	db := &testDB{Name: "test"}
	c.Provide("db", db)

	got, ok := c.Get("db")
	if !ok {
		t.Fatal("expected to find 'db'")
	}

	gotDB, ok := got.(*testDB)
	if !ok {
		t.Fatal("expected *testDB type")
	}

	if gotDB.Name != "test" {
		t.Errorf("got %q, want %q", gotDB.Name, "test")
	}
}

func TestContainer_Has(t *testing.T) {
	c := New()
	c.Provide("exists", "value")

	if !c.Has("exists") {
		t.Error("Has should return true for 'exists'")
	}

	if c.Has("missing") {
		t.Error("Has should return false for 'missing'")
	}
}

func TestContainer_Inject(t *testing.T) {
	c := New()
	db := &testDB{Name: "injected"}
	c.Provide("db", db)
	c.Provide("logger", "stdout")

	svc := &testService{}
	c.Inject(svc)

	if svc.DB != db {
		t.Errorf("DB not injected correctly")
	}

	if svc.Logger != "stdout" {
		t.Errorf("Logger: got %q, want %q", svc.Logger, "stdout")
	}
}

func TestContainer_InjectFallbackToFieldName(t *testing.T) {
	type noTagService struct {
		DB *testDB
	}

	c := New()
	db := &testDB{Name: "fallback"}
	c.Provide("DB", db)

	svc := &noTagService{}
	c.Inject(svc)

	if svc.DB != db {
		t.Error("DB not injected via field name fallback")
	}
}

func TestContainer_Clone(t *testing.T) {
	c := New()
	c.Provide("key", "value")

	clone := c.Clone()
	clone.Provide("new", "added")

	if !clone.Has("key") {
		t.Error("clone should have 'key'")
	}

	if !clone.Has("new") {
		t.Error("clone should have 'new'")
	}

	if c.Has("new") {
		t.Error("original should not have 'new'")
	}
}

func TestContainer_Keys(t *testing.T) {
	c := New()
	c.Provide("a", 1)
	c.Provide("b", 2)

	keys := c.Keys()
	if len(keys) != 2 {
		t.Errorf("got %d keys, want 2", len(keys))
	}
}

func TestContainer_MustGet_Panic(t *testing.T) {
	c := New()

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGet should panic for missing key")
		}
	}()

	c.MustGet("missing")
}

func TestResolve(t *testing.T) {
	c := New()
	c.Provide("num", 42)

	got, ok := Resolve[int](c, "num")
	if !ok {
		t.Fatal("expected to resolve 'num'")
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
}

func TestResolve_TypeMismatch(t *testing.T) {
	c := New()
	c.Provide("num", 42)

	_, ok := Resolve[string](c, "num")
	if ok {
		t.Error("should return false for type mismatch")
	}
}

func TestMustResolve_Panic(t *testing.T) {
	c := New()

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustResolve should panic for missing key")
		}
	}()

	MustResolve[int](c, "missing")
}
