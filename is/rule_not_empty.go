package is

import (
	"context"
	"reflect"
	"github.com/alexisvisco/valid/ishelper"
)

// NotEmpty returns a Rule that reports a violation when len(value) == 0.
//
// Accepted types: string, array, slice, map.
// Nil and unsupported types produce ViolationNotEmpty.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
func NotEmpty(_ context.Context, value any) *Violation {
	resolved, skip := ishelper.ExtractOptional(value)
	if skip {
		return nil
	}
	if resolved == nil {
		return &Violation{
			Code:    ViolationNotEmpty,
			Message: formatMessage(ViolationNotEmpty, nil),
		}
	}

	rv := reflect.ValueOf(resolved)
	switch rv.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		if rv.Len() == 0 {
			return &Violation{
				Code:    ViolationNotEmpty,
				Message: formatMessage(ViolationNotEmpty, nil),
			}
		}
		return nil
	default:
		return &Violation{
			Code:    ViolationNotEmpty,
			Message: formatMessage(ViolationNotEmpty, nil),
		}
	}
}
