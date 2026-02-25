package is

import (
	"context"
	"github.com/alexisvisco/valid/ishelper"
)

// GreaterThanOrEqual returns a Rule that reports a violation when value is < limit.
//
// Accepted types: all types in ishelper.Number (signed/unsigned integers, float32, float64).
// Unsupported/non-numeric values produce ViolationGTE.
//
// Optional behaviour: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
func GreaterThanOrEqual[T ishelper.Number](limit T) Rule {
	boundary, ok := ishelper.ToRat(limit)
	if !ok {
		panic("is.GreaterThanOrEqual: invalid limit value")
	}

	return func(_ context.Context, value any) *Violation {
		resolved, skip := ishelper.ExtractOptional(value)
		if skip {
			return nil
		}
		n, ok := ishelper.ToRat(resolved)
		if !ok || n.Cmp(boundary) < 0 {
			return &Violation{
				Code:    ViolationGTE,
				Message: formatMessage(ViolationGTE, map[string]any{"value": limit}),
			}
		}
		return nil
	}
}
