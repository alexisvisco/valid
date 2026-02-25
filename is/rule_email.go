package is

import (
	"context"
	"net/mail"
	"github.com/alexisvisco/valid/ishelper"
)

// Email reports a violation when value is not a valid email address string.
//
// Accepted type: string.
// Unsupported types and malformed email text produce ViolationEmail.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
func Email(_ context.Context, value any) *Violation {
	resolved, skip := ishelper.ExtractOptional(value)
	if skip {
		return nil
	}
	s, ok := resolved.(string)
	if !ok {
		return &Violation{Code: ViolationEmail, Message: formatMessage(ViolationEmail, nil)}
	}
	addr, err := mail.ParseAddress(s)
	if err != nil || addr.Address != s {
		return &Violation{Code: ViolationEmail, Message: formatMessage(ViolationEmail, nil)}
	}
	return nil
}
