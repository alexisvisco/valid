package is

import (
	"context"
	"regexp"
	"github.com/alexisvisco/valid/ishelper"
)

var numericRegex = regexp.MustCompile(`^[-+]?[0-9]+(?:\.[0-9]+)?$`)

// Numeric is a Rule that reports a violation when value is not numeric text.
//
// Accepted type: string.
// Unsupported types and non-matching text produce ViolationNumeric.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
var Numeric Rule = func(_ context.Context, value any) *Violation {
	resolved, skip := ishelper.ExtractOptional(value)
	if skip {
		return nil
	}
	s, ok := resolved.(string)
	if !ok || !numericRegex.MatchString(s) {
		return &Violation{Code: ViolationNumeric, Message: formatMessage(ViolationNumeric, nil)}
	}
	return nil
}
