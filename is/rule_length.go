package is

import (
	"context"
	"reflect"
	"valid/ishelper"
)

// Length returns a Rule that reports a violation when the length of the value
// falls outside the inclusive range [min, max].
//
// Accepted types: string, slice, array, map.
// Nil and unsupported types produce ViolationLength.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
func Length(min, max int) Rule {
	return func(_ context.Context, value any) *Violation {
		resolved, skip := ishelper.ExtractOptional(value)
		if skip {
			return nil
		}

		if resolved == nil {
			return &Violation{
				Code:    ViolationLength,
				Message: formatMessage(ViolationLength, map[string]any{"min": min, "max": max}),
			}
		}

		rv := reflect.ValueOf(resolved)
		switch rv.Kind() {
		case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
			l := rv.Len()
			if l < min || l > max {
				return &Violation{
					Code:    ViolationLength,
					Message: formatMessage(ViolationLength, map[string]any{"min": min, "max": max}),
				}
			}
			return nil
		default:
			return &Violation{
				Code:    ViolationLength,
				Message: formatMessage(ViolationLength, map[string]any{"min": min, "max": max}),
			}
		}
	}
}
