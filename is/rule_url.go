package is

import (
	"context"
	"net/url"
	"valid/ishelper"
)

// URL reports a violation when value is not an absolute URL string.
//
// Accepted type: string.
// Unsupported types, relative URLs, and malformed URL text produce ViolationURL.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
func URL(_ context.Context, value any) *Violation {
	resolved, skip := ishelper.ExtractOptional(value)
	if skip {
		return nil
	}
	s, ok := resolved.(string)
	if !ok {
		return &Violation{Code: ViolationURL, Message: formatMessage(ViolationURL, nil)}
	}
	u, err := url.ParseRequestURI(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return &Violation{Code: ViolationURL, Message: formatMessage(ViolationURL, nil)}
	}
	return nil
}
