package is

import (
	"context"
	"valid/ishelper"
)

// Min returns a Rule that reports a violation when the numeric value is strictly
// less than min.
//
// Accepted types: ishelper.Number (signed/unsigned integers, float32, float64).
// Unsupported/non-numeric values produce ViolationMin.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
func Min[T ishelper.Number](min T) Rule {
	limit, ok := ishelper.ToRat(min)
	if !ok {
		panic("is.Min: invalid min value")
	}

	return func(_ context.Context, value any) *Violation {
		resolved, skip := ishelper.ExtractOptional(value)
		if skip {
			return nil
		}

		n, ok := ishelper.ToRat(resolved)
		if !ok || n.Cmp(limit) < 0 {
			return &Violation{
				Code:    ViolationMin,
				Message: formatMessage(ViolationMin, map[string]any{"min": min}),
			}
		}
		return nil
	}
}
