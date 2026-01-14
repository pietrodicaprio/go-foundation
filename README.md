# go-foundation

A collection of foundational Go primitives for building declarative, struct-tag-driven libraries.

This SDK provides shared building blocks used across my ecosystem, including:
- **Struct Tag Parsing**: Generic `key:value` tag parser
- **Dependency Injection**: Minimal, type-safe DI container
- **Adapter Pattern**: Pluggable backend registry
- **Hook Discovery**: Automatic lifecycle hook detection

## Installation

```bash
go get github.com/mirkobrombin/go-foundation
```

## Modules

### `pkg/tags` - Struct Tag Parser

Generic parser for struct tags with `key:value` syntax.

```go
import "github.com/mirkobrombin/go-foundation/pkg/tags"

p := tags.NewParser("guard", tags.WithPairDelimiter(";"))
result := p.Parse("role:owner; read:admin,user")
// result["role"] = ["owner"]
// result["read"] = ["admin", "user"]
```

### `pkg/di` - Dependency Injection

Minimal DI container with generics support.

```go
import "github.com/mirkobrombin/go-foundation/pkg/di"

c := di.New()
c.Provide("db", myDB)

db := di.Get[*sql.DB](c, "db")
```

### `pkg/adapters` - Pluggable Backends

Generic registry for swappable adapters/backends.

```go
import "github.com/mirkobrombin/go-foundation/pkg/adapters"

r := adapters.NewRegistry[Transport]()
r.Register("http", httpTransport)
r.Register("grpc", grpcTransport)
r.SetDefault("http")

t := r.Default()
```

### `pkg/hooks` - Lifecycle Hooks

Automatic discovery of lifecycle methods via reflection.

```go
import "github.com/mirkobrombin/go-foundation/pkg/hooks"

d := hooks.NewDiscovery()
methods := d.Discover(myStruct, "OnEnter")
// Returns map of "OnEnterPaid", "OnEnterCancelled", etc.
```

## Why go-foundation?

This library consolidates patterns that were duplicated across multiple of my projects, I just
thought it would be a good idea to have a shared library for these primitives.

## Migration from go-struct-flags

`go-struct-flags` is now deprecated. Replace:

```go
// Before
import "github.com/mirkobrombin/go-struct-flags/v2/pkg/binder"

// After
import "github.com/mirkobrombin/go-foundation/pkg/tags"
```

as for now, they are fully compatible and can be used interchangeably.

## License

MIT License. See [LICENSE](LICENSE) for details.
