package is

import (
	"context"
	"strings"
	"github.com/alexisvisco/valid/ishelper"
)

// HasPrefix returns a Rule that reports a violation when the string value does
// not start with the given prefix.
//
// Accepted type: string.
// Unsupported/non-string values produce ViolationHasPrefix.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
func HasPrefix(prefix string) Rule {
	return func(_ context.Context, value any) *Violation {
		resolved, skip := ishelper.ExtractOptional(value)
		if skip {
			return nil
		}

		s, ok := resolved.(string)
		if !ok || !strings.HasPrefix(s, prefix) {
			return &Violation{
				Code:    ViolationHasPrefix,
				Message: formatMessage(ViolationHasPrefix, map[string]any{"prefix": prefix}),
			}
		}
		return nil
	}
}
