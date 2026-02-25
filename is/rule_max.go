package is

import (
	"context"
	"github.com/alexisvisco/valid/ishelper"
)

// Max returns a Rule that reports a violation when the numeric value is strictly
// greater than max.
//
// Accepted types: ishelper.Number (signed/unsigned integers, float32, float64).
// Unsupported/non-numeric values produce ViolationMax.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
func Max[T ishelper.Number](max T) Rule {
	limit, ok := ishelper.ToRat(max)
	if !ok {
		panic("is.Max: invalid max value")
	}

	return func(_ context.Context, value any) *Violation {
		resolved, skip := ishelper.ExtractOptional(value)
		if skip {
			return nil
		}

		n, ok := ishelper.ToRat(resolved)
		if !ok || n.Cmp(limit) > 0 {
			return &Violation{
				Code:    ViolationMax,
				Message: formatMessage(ViolationMax, map[string]any{"max": max}),
			}
		}
		return nil
	}
}
