# go-cqrs

[![CI](https://github.com/jassus213/go-cqrs/actions/workflows/ci.yml/badge.svg)](https://github.com/jassus213/go-cqrs/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/jassus213/go-cqrs/badge.svg?branch=main&v=2)](https://coveralls.io/github/jassus213/go-cqrs?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/jassus213/go-cqrs.svg)](https://pkg.go.dev/github.com/jassus213/go-cqrs)
[![Go Report Card](https://goreportcard.com/badge/github.com/jassus213/go-cqrs?v=2)](https://goreportcard.com/report/github.com/jassus213/go-cqrs)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Generic, composable use-case pipelines for Go. Zero dependencies.

```
go get github.com/jassus213/go-cqrs
```

## What it does

Wraps your business logic in a `UseCase[Req, Res]` function and lets you stack
cross-cutting concerns (logging, validation, recovery, …) as decorators — like
HTTP middleware, but for any operation.

```
Recovery → Logging → Validation → YourHandler
```

## Quick start

```go
package main

import (
    "context"
    cqrs "github.com/jassus213/go-cqrs"
)

type GetUserQuery struct{ ID int64 }
type User struct{ ID int64; Name string }

func getUser(_ context.Context, req GetUserQuery) (User, error) {
    return User{ID: req.ID, Name: "Alice"}, nil
}

func main() {
    uc := cqrs.NewDefaultBuilder(logger, getUser).Build()

    user, err := uc(context.Background(), GetUserQuery{ID: 1})
    // ...
}
```

## Core types

| Type | Description |
|------|-------------|
| `UseCase[Req, Res]` | `func(ctx, req) (res, error)` — your handler |
| `UseCaseDecorator[Req, Res]` | Wraps a `UseCase` and returns a new one |
| `UseCaseBuilder[Req, Res]` | Chains decorators and builds the final pipeline |
| `None` | Alias for `struct{}` — use as `Res` for commands |
| `Logger` | Minimal interface: `Info`, `Debug`, `Error` |
| `Validator` | Optional interface: request types implement `Validate() error` |

## Built-in decorators

| Decorator | What it does |
|-----------|-------------|
| `Recovery(logger)` | Catches panics, logs them, returns an error |
| `Logging(logger)` | Logs operation name, elapsed time, errors |
| `Validation()` | Calls `Validate()` on requests that implement `Validator` |

### NewDefaultBuilder

Pre-configured pipeline: **Recovery → Logging → Validation → handler**.

```go
uc := cqrs.NewDefaultBuilder(logger, handler).Build()
```

### Custom pipeline

```go
uc := cqrs.NewBuilder(handler).
    Use(
        cqrs.Recovery[Req, Res](logger),
        MyCustomDecorator[Req, Res](),
    ).
    Build()
```

## Bring your own Logger

Implement the `Logger` interface with any logging library:

```go
// slog adapter (3 lines)
type SlogAdapter struct{ L *slog.Logger }

func (a SlogAdapter) Info(ctx context.Context, msg string, args ...any)  { a.L.InfoContext(ctx, msg, args...) }
func (a SlogAdapter) Debug(ctx context.Context, msg string, args ...any) { a.L.DebugContext(ctx, msg, args...) }
func (a SlogAdapter) Error(ctx context.Context, msg string, args ...any) { a.L.ErrorContext(ctx, msg, args...) }
```

## Writing custom decorators

A decorator is just a function that wraps a `UseCase`:

```go
func Timeout[Req, Res any](d time.Duration) cqrs.UseCaseDecorator[Req, Res] {
    return func(next cqrs.UseCase[Req, Res]) cqrs.UseCase[Req, Res] {
        return func(ctx context.Context, req Req) (Res, error) {
            ctx, cancel := context.WithTimeout(ctx, d)
            defer cancel()
            return next(ctx, req)
        }
    }
}
```

## Commands (no return value)

Use `cqrs.None` as the response type:

```go
var deleteUser cqrs.UseCase[DeleteCmd, cqrs.None] = func(ctx context.Context, req DeleteCmd) (cqrs.None, error) {
    return cqrs.None{}, repo.Delete(ctx, req.ID)
}
```

## Dependency Injection (uber/fx)

go-cqrs works naturally with DI containers. Here's a production pattern with [uber/fx](https://github.com/uber-go/fx):

### 1. Define your use-case constructor

```go
package bll

import (
    "context"
    "fmt"

    cqrs "github.com/jassus213/go-cqrs"
    "go.uber.org/fx"
)

type GetAccountQuery struct {
    ID int64
}

// Dependencies are injected by fx via fx.In
type GetAccountDeps struct {
    fx.In

    Repository AccountRepository
    Logger     cqrs.Logger
}

func NewGetAccountUseCase(deps GetAccountDeps) cqrs.UseCase[GetAccountQuery, Account] {
    handler := func(ctx context.Context, req GetAccountQuery) (Account, error) {
        account, err := deps.Repository.FindByID(ctx, req.ID)
        if err != nil {
            return Account{}, fmt.Errorf("get account: %w", err)
        }
        return account, nil
    }

    return cqrs.NewDefaultBuilder(deps.Logger, handler).Build()
}
```

### 2. Register in an fx module

```go
package account

var Module = fx.Module("account",
    fx.Provide(
        fx.Annotate(
            postgres.NewAccountRepository,
            fx.As(new(AccountRepository)),
        ),
    ),
    fx.Provide(bll.NewGetAccountUseCase),
    fx.Provide(v1.NewAccountHandler),
)
```

### 3. Inject the built use-case into your handler

```go
package v1

type AccountHandler struct {
    getAccount cqrs.UseCase[bll.GetAccountQuery, domain.Account]
}

func NewAccountHandler(getAccount cqrs.UseCase[bll.GetAccountQuery, domain.Account]) *AccountHandler {
    return &AccountHandler{getAccount: getAccount}
}

func (h *AccountHandler) Get(c *gin.Context) {
    account, err := h.getAccount(c.Request.Context(), bll.GetAccountQuery{ID: 1})
    // ...
}
```

The key insight: `fx.In` tells fx to resolve each struct field from the container. The constructor receives fully-injected deps, builds the pipeline once at startup, and returns a ready-to-call `UseCase` function. No runtime overhead.

## Examples

See [`examples/basic/main.go`](examples/basic/main.go) for a runnable demo.

## License

MIT
