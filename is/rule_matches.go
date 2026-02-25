package is

import (
	"context"
	"regexp"
	"github.com/alexisvisco/valid/ishelper"
)

// Matches returns a Rule that reports a violation when value does not match pattern.
//
// Accepted type: string.
// Unsupported types and non-matching text produce ViolationMatches.
//
// The pattern must be a valid Go regexp. Invalid patterns panic at rule construction time.
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
func Matches(pattern string) Rule {
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic("is.Matches: invalid pattern")
	}
	return func(_ context.Context, value any) *Violation {
		resolved, skip := ishelper.ExtractOptional(value)
		if skip {
			return nil
		}
		s, ok := resolved.(string)
		if !ok || !re.MatchString(s) {
			return &Violation{
				Code:    ViolationMatches,
				Message: formatMessage(ViolationMatches, map[string]any{"pattern": pattern}),
			}
		}
		return nil
	}
}
