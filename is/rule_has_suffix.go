package is

import (
	"context"
	"strings"
	"github.com/alexisvisco/valid/ishelper"
)

// HasSuffix returns a Rule that reports a violation when the string value does
// not end with the given suffix.
//
// Accepted type: string.
// Unsupported/non-string values produce ViolationHasSuffix.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
func HasSuffix(suffix string) Rule {
	return func(_ context.Context, value any) *Violation {
		resolved, skip := ishelper.ExtractOptional(value)
		if skip {
			return nil
		}

		s, ok := resolved.(string)
		if !ok || !strings.HasSuffix(s, suffix) {
			return &Violation{
				Code:    ViolationHasSuffix,
				Message: formatMessage(ViolationHasSuffix, map[string]any{"suffix": suffix}),
			}
		}
		return nil
	}
}
