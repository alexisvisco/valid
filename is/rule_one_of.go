package is

import (
	"context"
	"fmt"
	"strings"
	"github.com/alexisvisco/valid/ishelper"
)

// OneOf returns a Rule that reports a violation when value is not one of allowed.
//
// Accepted type: T, where T is any comparable type.
// Type mismatches and values outside the allowed set produce ViolationOneOf.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value against the allowed list.
func OneOf[T comparable](allowed ...T) Rule {
	return func(_ context.Context, value any) *Violation {
		resolved, skip := ishelper.ExtractOptional(value)
		if skip {
			return nil
		}

		v, ok := resolved.(T)
		if ok {
			for _, a := range allowed {
				if v == a {
					return nil
				}
			}
		}

		parts := make([]string, len(allowed))
		for i, a := range allowed {
			parts[i] = fmt.Sprintf("%v", a)
		}
		return &Violation{
			Code:    ViolationOneOf,
			Message: formatMessage(ViolationOneOf, map[string]any{"values": strings.Join(parts, ", ")}),
		}
	}
}
