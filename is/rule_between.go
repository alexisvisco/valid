package is

import (
	"context"
	"valid/ishelper"
)

// Between returns a Rule that reports a violation when value is outside [min, max].
//
// Accepted types: ishelper.Number (signed/unsigned integers, float32, float64).
// Unsupported/non-numeric values produce ViolationBetween.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
func Between[T ishelper.Number](min, max T) Rule {
	minLimit, ok := ishelper.ToRat(min)
	if !ok {
		panic("is.Between: invalid min value")
	}
	maxLimit, ok := ishelper.ToRat(max)
	if !ok {
		panic("is.Between: invalid max value")
	}

	return func(_ context.Context, value any) *Violation {
		resolved, skip := ishelper.ExtractOptional(value)
		if skip {
			return nil
		}
		n, ok := ishelper.ToRat(resolved)
		if !ok || n.Cmp(minLimit) < 0 || n.Cmp(maxLimit) > 0 {
			return &Violation{
				Code:    ViolationBetween,
				Message: formatMessage(ViolationBetween, map[string]any{"min": min, "max": max}),
			}
		}
		return nil
	}
}
