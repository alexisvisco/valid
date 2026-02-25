package ishelper

import "reflect"

// Optional is the interface that optional values must implement to be detected by is rules.
type Optional interface {
	IsSome() bool
	IsNone() bool
}

// ExtractOptional resolves an optional value for constraint rules (Min, Max, Length).
// Returns (nil, true) if the value is None — the rule should return nil (no violation).
// Returns (unwrapped, false) if the value is Some — the rule should validate the inner value.
// Returns (value, false) if the value is not optional — the rule should validate as-is.
func ExtractOptional(value any) (any, bool) {
	opt, ok := value.(Optional)
	if !ok {
		return value, false
	}
	if opt.IsNone() {
		return nil, true
	}
	return unwrap(value), false
}

// unwrap calls Unwrap() via reflection since Unwrap() is generic and cannot appear in a non-generic interface.
func unwrap(value any) any {
	rv := reflect.ValueOf(value)
	if m := rv.MethodByName("Unwrap"); m.IsValid() {
		if results := m.Call(nil); len(results) > 0 {
			return results[0].Interface()
		}
	}
	return value
}

type someVal[T any] struct{ v T }

func (s someVal[T]) IsSome() bool { return true }
func (s someVal[T]) IsNone() bool { return false }
func (s someVal[T]) Unwrap() T    { return s.v }

type noneVal[T any] struct{}

func (noneVal[T]) IsSome() bool { return false }
func (noneVal[T]) IsNone() bool { return true }
func (noneVal[T]) Unwrap() T    { var zero T; return zero }

// Some creates an optional value that is present.
func Some[T any](v T) someVal[T] { return someVal[T]{v: v} }

// None creates an optional value that is absent.
func None[T any]() noneVal[T] { return noneVal[T]{} }
