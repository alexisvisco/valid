package is

import (
	"context"
	"regexp"
	"github.com/alexisvisco/valid/ishelper"
)

var alphaNumericRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

// Alphanumeric is a Rule that reports a violation when value is not alphanumeric text.
//
// Accepted type: string.
// Unsupported types and non-matching text produce ViolationAlphaNum.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
var Alphanumeric Rule = func(_ context.Context, value any) *Violation {
	resolved, skip := ishelper.ExtractOptional(value)
	if skip {
		return nil
	}
	s, ok := resolved.(string)
	if !ok || !alphaNumericRegex.MatchString(s) {
		return &Violation{Code: ViolationAlphaNum, Message: formatMessage(ViolationAlphaNum, nil)}
	}
	return nil
}
