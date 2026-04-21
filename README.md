# go-cqrs

[![CI](https://github.com/jassus213/go-cqrs/actions/workflows/ci.yml/badge.svg)](https://github.com/jassus213/go-cqrs/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/jassus213/go-cqrs/badge.svg?branch=main&v=2)](https://coveralls.io/github/jassus213/go-cqrs?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/jassus213/go-cqrs.svg)](https://pkg.go.dev/github.com/jassus213/go-cqrs)
[![Go Report Card](https://goreportcard.com/badge/github.com/jassus213/go-cqrs?v=2)](https://goreportcard.com/report/github.com/jassus213/go-cqrs)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Generic, composable CQRS pipelines for Go. Zero dependencies.

```
go get github.com/jassus213/go-cqrs
```

## What it does

Wraps your business logic in typed `QueryHandler` or `CommandHandler` functions and lets you
stack cross-cutting concerns (logging, validation, recovery, …) as decorators — like HTTP
middleware, but for any operation. Queries (reads) and commands (writes) are explicit, separate
types for clarity.

```
QueryRecovery → QueryLogging → QueryValidation → YourQueryHandler
CommandRecovery → CommandLogging → CommandValidation → YourCommandHandler
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
    h := cqrs.NewDefaultQueryBuilder(logger, cqrs.QueryHandler[GetUserQuery, User](getUser)).Build()

    user, err := h(context.Background(), GetUserQuery{ID: 1})
    // ...
}
```

## Core types

| Type | Description |
|------|-------------|
| `QueryHandler[Req, Res]` | `func(ctx, req) (res, error)` — read operation |
| `CommandHandler[Req, Res]` | `func(ctx, req) (res, error)` — write operation |
| `QueryDecorator[Req, Res]` | Wraps a `QueryHandler` and returns a new one |
| `CommandDecorator[Req, Res]` | Wraps a `CommandHandler` and returns a new one |
| `None` | Alias for `struct{}` — use as `Res` for commands that return nothing |
| `Logger` | Minimal interface: `Info`, `Warn`, `Debug`, `Error` |
| `Validator` | Optional interface: request types implement `Validate() error` |

## Built-in decorators

| Decorator | What it does |
|-----------|-------------|
| `QueryRecovery(logger)` / `CommandRecovery(logger)` | Catches panics, logs them, returns an error |
| `QueryLogging(logger)` / `CommandLogging(logger)` | Logs operation name, elapsed time, errors |
| `QueryValidation()` / `CommandValidation()` | Calls `Validate()` on requests that implement `Validator` |

### Default pipelines

Pre-configured stacks: **Recovery → Logging → Validation → handler**.

```go
// Query
h := cqrs.NewDefaultQueryBuilder(logger, queryHandler).Build()

// Command
h := cqrs.NewDefaultCommandBuilder(logger, commandHandler).Build()
```

### Custom pipeline

```go
h := cqrs.NewQueryBuilder(handler).
    Use(
        cqrs.QueryRecovery[Req, Res](logger),
        MyCustomDecorator[Req, Res](),
    ).
    Build()
```

## Bring your own Logger

Implement the `Logger` interface with any logging library:

```go
// slog adapter
type SlogAdapter struct{ L *slog.Logger }

func (a SlogAdapter) Info(ctx context.Context, msg string, args ...any)  { a.L.InfoContext(ctx, msg, args...) }
func (a SlogAdapter) Warn(ctx context.Context, msg string, args ...any)  { a.L.WarnContext(ctx, msg, args...) }
func (a SlogAdapter) Debug(ctx context.Context, msg string, args ...any) { a.L.DebugContext(ctx, msg, args...) }
func (a SlogAdapter) Error(ctx context.Context, msg string, args ...any) { a.L.ErrorContext(ctx, msg, args...) }
```

## Writing custom decorators

A decorator is a function that wraps a handler:

```go
func Timeout[Req, Res any](d time.Duration) cqrs.QueryDecorator[Req, Res] {
    return func(next cqrs.QueryHandler[Req, Res]) cqrs.QueryHandler[Req, Res] {
        return func(ctx context.Context, req Req) (Res, error) {
            ctx, cancel := context.WithTimeout(ctx, d)
            defer cancel()
            return next(ctx, req)
        }
    }
}
```

The same pattern applies to `CommandDecorator` / `CommandHandler`.

## Commands (no return value)

Use `cqrs.None` as the response type:

```go
var deleteUser cqrs.CommandHandler[DeleteCmd, cqrs.None] = func(ctx context.Context, req DeleteCmd) (cqrs.None, error) {
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

type GetAccountDeps struct {
    fx.In

    Repository AccountRepository
    Logger     cqrs.Logger
}

func NewGetAccountHandler(deps GetAccountDeps) cqrs.QueryHandler[GetAccountQuery, Account] {
    handler := func(ctx context.Context, req GetAccountQuery) (Account, error) {
        account, err := deps.Repository.FindByID(ctx, req.ID)
        if err != nil {
            return Account{}, fmt.Errorf("get account: %w", err)
        }
        return account, nil
    }

    return cqrs.NewDefaultQueryBuilder(deps.Logger, handler).Build()
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
    fx.Provide(bll.NewGetAccountHandler),
    fx.Provide(v1.NewAccountHandler),
)
```

### 3. Inject into your HTTP handler

```go
package v1

type AccountHandler struct {
    getAccount cqrs.QueryHandler[bll.GetAccountQuery, domain.Account]
}

func NewAccountHandler(getAccount cqrs.QueryHandler[bll.GetAccountQuery, domain.Account]) *AccountHandler {
    return &AccountHandler{getAccount: getAccount}
}

func (h *AccountHandler) Get(c *gin.Context) {
    account, err := h.getAccount(c.Request.Context(), bll.GetAccountQuery{ID: 1})
    // ...
}
```

## Examples

See [`examples/basic/main.go`](examples/basic/main.go) for a runnable demo.

## License

MIT
