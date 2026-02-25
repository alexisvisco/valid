package is

import (
	"context"
	"reflect"
	"strings"
	"github.com/alexisvisco/valid/ishelper"
)

// Contains returns a Rule that reports a violation when the value does not
// contain elem. The behavior depends on the runtime type of the value:
//
//   - string value + string elem: violation if the string does not contain elem as substring.
//   - slice/array value: violation if none of the elements equals elem.
//
// Accepted types: string, slice, array (with comparable element type).
// Any other combination of types produces ViolationContains.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
func Contains[T comparable](elem T) Rule {
	return func(_ context.Context, value any) *Violation {
		resolved, skip := ishelper.ExtractOptional(value)
		if skip {
			return nil
		}

		violation := &Violation{
			Code:    ViolationContains,
			Message: formatMessage(ViolationContains, map[string]any{"value": elem}),
		}

		// String contains substring (only when elem is a string).
		if s, ok := resolved.(string); ok {
			if sub, ok := any(elem).(string); ok {
				if strings.Contains(s, sub) {
					return nil
				}
				return violation
			}
		}

		// Slice or array contains element.
		rv := reflect.ValueOf(resolved)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			for i := range rv.Len() {
				if item, ok := rv.Index(i).Interface().(T); ok && item == elem {
					return nil
				}
			}
			return violation
		default:
			return violation
		}
	}
}
