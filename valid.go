package valid

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"github.com/alexisvisco/valid/is"
)

// FieldError represents a single field validation error.
type FieldError struct {
	Path    string
	Code    string
	Message string
}

// Error is a collection of FieldErrors returned by Struct.
type Error struct {
	Fields []FieldError
}

func (e *Error) Error() string {
	parts := make([]string, len(e.Fields))
	for i, fe := range e.Fields {
		parts[i] = fmt.Sprintf("%s (%s)", fe.Path, fe.Code)
	}
	return "validation failed: " + strings.Join(parts, ", ")
}

// Rename returns a new *Error with field paths replaced according to mapping.
// "*" in mapping keys acts as a wildcard matching any single path segment (e.g. array indices).
// Fields without a match keep their original path. Returns nil if e is nil.
func (e *Error) Rename(mapping map[string]string) *Error {
	if e == nil {
		return nil
	}
	var fields []FieldError
	for _, fe := range e.Fields {
		target, _ := matchAndRename(fe.Path, mapping)
		if target == "" {
			target = fe.Path
		}
		fields = append(fields, FieldError{
			Path:    target,
			Code:    fe.Code,
			Message: fe.Message,
		})
	}
	if len(fields) == 0 {
		return nil
	}
	return &Error{Fields: fields}
}

// FieldGroup is a lazy field validation group. It is created by Field, Slice,
// Each, and Nested, and evaluated by Struct with a context.
type FieldGroup func(ctx context.Context) []FieldError

// Validatable is implemented by types that can validate themselves.
// The context is forwarded from the enclosing Struct call.
type Validatable interface {
	Valid(ctx context.Context) error
}

// Field returns a FieldGroup that evaluates the given rules against value when
// called by Struct. Rules are short-circuited: the first violation stops evaluation.
func Field(path string, value any, rules ...is.Rule) FieldGroup {
	return func(ctx context.Context) []FieldError {
		for _, rule := range rules {
			if v := rule(ctx, value); v != nil {
				return []FieldError{{
					Path:    path,
					Code:    string(v.Code),
					Message: v.Message,
				}}
			}
		}
		return nil
	}
}

// Struct evaluates all groups with ctx and aggregates their FieldErrors into a
// single *Error. Returns nil if no errors are found.
//
// Path deduplication: once a path X has an error (from any group), subsequent
// groups' errors for X or any child path X.* are skipped. This prevents
// cascading errors when a field-level check (e.g. Required) is paired with a
// nested check (e.g. Nested) for the same path.
func Struct(ctx context.Context, groups ...FieldGroup) error {
	seenPaths := map[string]bool{}
	var all []FieldError
	for _, g := range groups {
		for _, e := range g(ctx) {
			if !hasFailedAncestor(e.Path, seenPaths) {
				all = append(all, e)
				seenPaths[e.Path] = true
			}
		}
	}
	if len(all) == 0 {
		return nil
	}
	return &Error{Fields: all}
}

// hasFailedAncestor reports whether any previously-seen path is an ancestor of
// (or equal to) path. E.g. "Payment" is an ancestor of "Payment.Method".
func hasFailedAncestor(path string, seen map[string]bool) bool {
	for p := range seen {
		if path == p || strings.HasPrefix(path, p+".") {
			return true
		}
	}
	return false
}

// Nested delegates validation to a value that implements Validatable.
// It accepts either a single Validatable or a slice of Validatable (detected at runtime).
// Nil pointers and nil slices produce no errors.
// Field errors from nested validation are prefixed with path.
func Nested(path string, v any) FieldGroup {
	return func(ctx context.Context) []FieldError {
		if v == nil {
			return nil
		}
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return nil
			}
			rv = rv.Elem()
		}
		if rv.Kind() == reflect.Slice {
			var errs []FieldError
			for i := 0; i < rv.Len(); i++ {
				itemPath := fmt.Sprintf("%s.%d", path, i)
				errs = append(errs, nestedOne(itemPath, rv.Index(i).Interface())(ctx)...)
			}
			return errs
		}
		return nestedOne(path, v)(ctx)
	}
}

func nestedOne(path string, v any) FieldGroup {
	return func(ctx context.Context) []FieldError {
		if v == nil {
			return nil
		}
		if rv := reflect.ValueOf(v); rv.Kind() == reflect.Ptr && rv.IsNil() {
			return nil
		}
		val, ok := v.(Validatable)
		if !ok {
			return nil
		}
		err := val.Valid(ctx)
		if err == nil {
			return nil
		}
		if ve := As(err); ve != nil {
			fields := make([]FieldError, len(ve.Fields))
			for i, fe := range ve.Fields {
				fields[i] = FieldError{
					Path:    path + "." + fe.Path,
					Code:    fe.Code,
					Message: fe.Message,
				}
			}
			return fields
		}
		return []FieldError{{Path: path, Code: "invalid"}}
	}
}

// Slice validates each element of items using fn and returns a FieldGroup.
// fn receives the context, the index, and the element; its FieldErrors are
// prefixed as "path.i.*" in the result.
func Slice[T any](path string, items []T, fn func(ctx context.Context, i int, item T) error) FieldGroup {
	return func(ctx context.Context) []FieldError {
		var errs []FieldError
		for i, item := range items {
			err := fn(ctx, i, item)
			if err == nil {
				continue
			}
			if ve := As(err); ve != nil {
				for _, fe := range ve.Fields {
					errs = append(errs, FieldError{
						Path:    fmt.Sprintf("%s.%d.%s", path, i, fe.Path),
						Code:    fe.Code,
						Message: fe.Message,
					})
				}
			} else {
				errs = append(errs, FieldError{
					Path: fmt.Sprintf("%s.%d", path, i),
					Code: "invalid",
				})
			}
		}
		return errs
	}
}

// Each validates each element of items against rules and returns a FieldGroup.
// Rules are short-circuited per element. Violations are reported as "path.i".
func Each[T any](path string, items []T, rules ...is.Rule) FieldGroup {
	return func(ctx context.Context) []FieldError {
		var errs []FieldError
		for i, item := range items {
			for _, rule := range rules {
				if v := rule(ctx, item); v != nil {
					errs = append(errs, FieldError{
						Path:    fmt.Sprintf("%s.%d", path, i),
						Code:    string(v.Code),
						Message: v.Message,
					})
					break
				}
			}
		}
		return errs
	}
}

// As returns *Error if err is (or wraps) a *Error. Returns nil otherwise.
func As(err error) *Error {
	var ve *Error
	if errors.As(err, &ve) {
		return ve
	}
	return nil
}

func matchAndRename(fieldPath string, mapping map[string]string) (string, bool) {
	fieldSegs := strings.Split(fieldPath, ".")
	for pattern, target := range mapping {
		patSegs := strings.Split(pattern, ".")
		if len(patSegs) != len(fieldSegs) {
			continue
		}
		var captures []string
		match := true
		for i, pat := range patSegs {
			if pat == "*" {
				captures = append(captures, fieldSegs[i])
			} else if pat != fieldSegs[i] {
				match = false
				break
			}
		}
		if !match {
			continue
		}
		result := target
		for _, captured := range captures {
			result = strings.Replace(result, "*", captured, 1)
		}
		return result, true
	}
	return "", false
}