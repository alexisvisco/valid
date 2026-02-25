package is

import (
	"context"
	"regexp"
	"github.com/alexisvisco/valid/ishelper"
)

var alphaRegex = regexp.MustCompile(`^[a-zA-Z]+$`)

// Alpha is a Rule that reports a violation when value is not alphabetic text.
//
// Accepted type: string.
// Unsupported types and non-matching text produce ViolationAlpha.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
var Alpha Rule = func(_ context.Context, value any) *Violation {
	resolved, skip := ishelper.ExtractOptional(value)
	if skip {
		return nil
	}
	s, ok := resolved.(string)
	if !ok || !alphaRegex.MatchString(s) {
		return &Violation{Code: ViolationAlpha, Message: formatMessage(ViolationAlpha, nil)}
	}
	return nil
}
