package is

import (
	"context"
	"reflect"
	"valid/ishelper"
)

// MinLength returns a Rule that reports a violation when len(value) < n.
//
// Accepted types: string, array, slice, map.
// Nil and unsupported types produce ViolationMinLength.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
func MinLength(n int) Rule {
	return func(_ context.Context, value any) *Violation {
		resolved, skip := ishelper.ExtractOptional(value)
		if skip {
			return nil
		}
		if resolved == nil {
			return &Violation{
				Code:    ViolationMinLength,
				Message: formatMessage(ViolationMinLength, map[string]any{"min": n}),
			}
		}

		rv := reflect.ValueOf(resolved)
		switch rv.Kind() {
		case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
			if rv.Len() < n {
				return &Violation{
					Code:    ViolationMinLength,
					Message: formatMessage(ViolationMinLength, map[string]any{"min": n}),
				}
			}
			return nil
		default:
			return &Violation{
				Code:    ViolationMinLength,
				Message: formatMessage(ViolationMinLength, map[string]any{"min": n}),
			}
		}
	}
}
