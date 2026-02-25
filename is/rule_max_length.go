package is

import (
	"context"
	"reflect"
	"valid/ishelper"
)

// MaxLength returns a Rule that reports a violation when len(value) > n.
//
// Accepted types: string, array, slice, map.
// Nil and unsupported types produce ViolationMaxLength.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
func MaxLength(n int) Rule {
	return func(_ context.Context, value any) *Violation {
		resolved, skip := ishelper.ExtractOptional(value)
		if skip {
			return nil
		}
		if resolved == nil {
			return &Violation{
				Code:    ViolationMaxLength,
				Message: formatMessage(ViolationMaxLength, map[string]any{"max": n}),
			}
		}

		rv := reflect.ValueOf(resolved)
		switch rv.Kind() {
		case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
			if rv.Len() > n {
				return &Violation{
					Code:    ViolationMaxLength,
					Message: formatMessage(ViolationMaxLength, map[string]any{"max": n}),
				}
			}
			return nil
		default:
			return &Violation{
				Code:    ViolationMaxLength,
				Message: formatMessage(ViolationMaxLength, map[string]any{"max": n}),
			}
		}
	}
}
