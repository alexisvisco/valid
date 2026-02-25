package is

import (
	"context"
	"math/big"
	"valid/ishelper"
)

// Positive is a Rule that reports a violation when value is <= 0.
//
// Accepted types: all numeric types supported by ishelper.ToRat.
// Unsupported/non-numeric values produce ViolationPositive.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
var Positive Rule = func(_ context.Context, value any) *Violation {
	resolved, skip := ishelper.ExtractOptional(value)
	if skip {
		return nil
	}
	n, ok := ishelper.ToRat(resolved)
	if !ok || n.Cmp(big.NewRat(0, 1)) <= 0 {
		return &Violation{
			Code:    ViolationPositive,
			Message: formatMessage(ViolationPositive, nil),
		}
	}
	return nil
}
