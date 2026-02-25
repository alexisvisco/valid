package is

import (
	"context"
	"math/big"
	"github.com/alexisvisco/valid/ishelper"
)

// NonNegative is a Rule that reports a violation when value is < 0.
//
// Accepted types: all numeric types supported by ishelper.ToRat.
// Unsupported/non-numeric values produce ViolationNonNeg.
//
// Optional behavior: None -> nil (absent field skips the constraint);
// Some(v) -> validates the unwrapped value.
var NonNegative Rule = func(_ context.Context, value any) *Violation {
	resolved, skip := ishelper.ExtractOptional(value)
	if skip {
		return nil
	}
	n, ok := ishelper.ToRat(resolved)
	if !ok || n.Cmp(big.NewRat(0, 1)) < 0 {
		return &Violation{
			Code:    ViolationNonNeg,
			Message: formatMessage(ViolationNonNeg, nil),
		}
	}
	return nil
}
