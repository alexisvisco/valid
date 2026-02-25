# valid

`valid` is a small, composable validation library for Go.

It gives you:
- field-level rules (`valid.Field` + `is.Rule`)
- aggregated validation errors (`*valid.Error`)
- nested struct and slice validation (`valid.Nested`, `valid.Slice`, `valid.Each`)
- path renaming for API-friendly error payloads (`(*valid.Error).Rename`)
- context-aware custom rules

## Install

This repository currently uses module path `valid`.

```bash
go get valid
```

Import packages:

```go
import (
    "valid"
    "valid/is"
)
```

## Quick start

```go
package main

import (
    "context"
    "fmt"

    "valid"
    "valid/is"
)

type CreateUserInput struct {
    Email string
    Name  string
    Age   int
}

func validateCreateUser(ctx context.Context, in CreateUserInput) error {
    return valid.Struct(ctx,
        valid.Field("Email", in.Email, is.Required, is.Email),
        valid.Field("Name", in.Name, is.Required, is.MinLength(2), is.MaxLength(50)),
        valid.Field("Age", in.Age, is.GreaterThanOrEqual(18)),
    )
}

func main() {
    err := validateCreateUser(context.Background(), CreateUserInput{})
    if err == nil {
        return
    }

    if ve := valid.As(err); ve != nil {
        for _, f := range ve.Fields {
            fmt.Printf("%s: %s (%s)\n", f.Path, f.Message, f.Code)
        }
    }
}
```

## Core concepts

### `valid.Field(path, value, rules...)`
Creates a lazy validation group for one field.
Rules are short-circuited: the first failing rule stops evaluation for that field.

### `valid.Struct(ctx, groups...)`
Runs all groups and returns:
- `nil` when everything passes
- `*valid.Error` when one or more fields fail

Path de-duplication is applied: once `X` fails, nested paths like `X.Y` are skipped.

### `*valid.Error` and `valid.As`
`valid.Struct` returns `error`; use `valid.As(err)` to safely extract `*valid.Error` (including wrapped errors).

## Nested validation

### `valid.Nested(path, v)` — delegate to `Validatable`

Implement the interface on your type:

```go
func (p PaymentParams) Valid(ctx context.Context) error {
    return valid.Struct(ctx,
        valid.Field("Method", p.Method, is.Required, is.OneOf("card", "bank_transfer")),
        valid.Field("TransactionID", p.TransactionID, is.Required, is.HasPrefix("txn_")),
    )
}
```

Then use `valid.Nested` for both single values and slices — it detects at runtime:

```go
valid.Struct(ctx,
    valid.Field("Payment", params.Payment, is.Required),
    valid.Nested("Payment", params.Payment),       // single Validatable
    valid.Nested("Discounts", params.Discounts),   // []Validatable
)
```

Field errors from nested validation are prefixed with the given path.
Nil pointers and nil slices produce no errors.

### `valid.Slice(path, items, fn)` — slice with custom validation

Use when each element needs its own `valid.Struct` call:

```go
valid.Slice("Items", params.Items, func(ctx context.Context, i int, item OrderItem) error {
    return valid.Struct(ctx,
        valid.Field("Name", item.Name, is.Required),
        valid.Field("Quantity", item.Quantity, is.Min(1)),
    )
})
```

If element `0` has a `Name` error, the output path is `Items.0.Name`.

### `valid.Each(path, items, rules...)` — per-element rules

Use when each element in a slice must satisfy the same rules:

```go
valid.Each("PaymentTypes", params.PaymentTypes, is.OneOf("card", "bank_transfer"))
```

Violations are reported with indexed paths: `PaymentTypes.2`.
Rules are short-circuited per element.

## Rename internal paths for public APIs

Use `(*valid.Error).Rename` to map internal field paths to response paths.
`*` matches exactly one path segment (e.g. an array index).
Fields without a mapping entry keep their original path.

```go
mapping := map[string]string{
    "BillingAddress.Street": "billingAddress.street",
    "Items.*.Name":          "items.*.name",
}

publicErr := valid.As(err).Rename(mapping)
```

## Custom rules with context

Rules have the signature:

```go
type Rule func(ctx context.Context, value any) *is.Violation
```

You can read values from `ctx` inside custom rules and compose them freely with built-in rules:

```go
var mustMatchTenant is.Rule = func(ctx context.Context, value any) *is.Violation {
    tenantID, _ := ctx.Value(tenantKey{}).(string)
    if value != tenantID {
        return &is.Violation{Code: "TENANT_MISMATCH"}
    }
    return nil
}
```

## Built-in rules (`valid/is`)

Each rule reports a violation code (e.g. `REQUIRED`, `MIN`, `EMAIL`) and a default message.

| Signature | Code | Accepted types | Description |
|---|---|---|---|
| `is.Required` | `REQUIRED` | any | Fails if value is nil, zero, or `None` |
| `is.NotEmpty` | `NOT_EMPTY` | string, slice, array, map | Fails if `len(value) == 0` |
| `is.Min(n T)` | `MIN` | integer, float | `value >= n` |
| `is.Max(n T)` | `MAX` | integer, float | `value <= n` |
| `is.Between(min, max T)` | `BETWEEN` | integer, float | `min <= value <= max` |
| `is.Positive` | `POSITIVE` | integer, float | `value > 0` |
| `is.NonNegative` | `NON_NEGATIVE` | integer, float | `value >= 0` |
| `is.GreaterThan(n T)` | `GT` | integer, float | `value > n` |
| `is.GreaterThanOrEqual(n T)` | `GTE` | integer, float | `value >= n` |
| `is.LessThan(n T)` | `LT` | integer, float | `value < n` |
| `is.LessThanOrEqual(n T)` | `LTE` | integer, float | `value <= n` |
| `is.Equal(target T)` | `EQ` | comparable | `value == target` |
| `is.MinLength(n int)` | `MIN_LENGTH` | string | `len(value) >= n` |
| `is.MaxLength(n int)` | `MAX_LENGTH` | string | `len(value) <= n` |
| `is.Length(min, max int)` | `LENGTH` | string | `min <= len(value) <= max` |
| `is.HasPrefix(s string)` | `HAS_PREFIX` | string | Value starts with `s` |
| `is.HasSuffix(s string)` | `HAS_SUFFIX` | string | Value ends with `s` |
| `is.Contains(elem T)` | `CONTAINS` | string (substring), slice, array | Value contains `elem` |
| `is.Matches(pattern string)` | `MATCHES` | string | Value matches regex pattern |
| `is.Alpha` | `ALPHA` | string | Only letters `[a-zA-Z]` |
| `is.Alphanumeric` | `ALPHANUMERIC` | string | Only letters and digits `[a-zA-Z0-9]` |
| `is.Numeric` | `NUMERIC` | string | Numeric text, e.g. `"123"`, `"-4.5"` |
| `is.Email` | `EMAIL` | string | Valid email address |
| `is.URL` | `URL` | string | Valid URL |
| `is.UUID` | `UUID` | string | Valid UUID (case-insensitive) |
| `is.OneOf(values ...T)` | `ONE_OF` | comparable | Value is one of the allowed values |

## Optional values

Most rules support optional values through `valid/ishelper`:

```go
import "valid/ishelper"

optName := ishelper.None[string]()
optQty  := ishelper.Some(10)
```

Behavior:
- `None` skips most constraint rules (field is absent)
- `is.Required` fails on `None`

## Run tests

```bash
go test ./...
```