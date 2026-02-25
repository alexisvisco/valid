package is

import (
	"context"
	"github.com/alexisvisco/valid/ishelper"
)

// Equal returns a Rule that reports a violation when value is not equal to target.
//
// Accepted types: any comparable type T.
// Values of a different runtime type, or values not equal to target, produce ViolationEQ.
//
// Optional behaviour: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
func Equal[T comparable](target T) Rule {
	return func(_ context.Context, value any) *Violation {
		resolved, skip := ishelper.ExtractOptional(value)
		if skip {
			return nil
		}

		v, ok := resolved.(T)
		if !ok || v != target {
			return &Violation{
				Code:    ViolationEQ,
				Message: formatMessage(ViolationEQ, map[string]any{"value": target}),
			}
		}
		return nil
	}
}
