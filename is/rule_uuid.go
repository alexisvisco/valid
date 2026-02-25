package is

import (
	"context"
	"regexp"
	"valid/ishelper"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// UUID is a Rule that reports a violation when value is not UUID text.
//
// Accepted type: string.
// Unsupported types and non-matching text produce ViolationUUID.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
var UUID Rule = func(_ context.Context, value any) *Violation {
	resolved, skip := ishelper.ExtractOptional(value)
	if skip {
		return nil
	}
	s, ok := resolved.(string)
	if !ok || !uuidRegex.MatchString(s) {
		return &Violation{Code: ViolationUUID, Message: formatMessage(ViolationUUID, nil)}
	}
	return nil
}
