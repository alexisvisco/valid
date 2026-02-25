package is

import (
	"context"
	"reflect"
	"github.com/alexisvisco/valid/ishelper"
)

// Required reports a violation when the value is absent or zero.
//
// For optional values (implementing ishelper.Optional):
//   - None -> ViolationRequired
//   - Some -> nil regardless of inner value (presence satisfies this rule)
//
// For non-optional values, accepted types are any.
// Nil, typed nil (pointer/slice/map/interface), and zero values produce ViolationRequired.
func Required(_ context.Context, value any) *Violation {
	if opt, ok := value.(ishelper.Optional); ok {
		if opt.IsNone() {
			return &Violation{
				Code:    ViolationRequired,
				Message: formatMessage(ViolationRequired, nil),
			}
		}
		return nil
	}

	if value == nil {
		return &Violation{
			Code:    ViolationRequired,
			Message: formatMessage(ViolationRequired, nil),
		}
	}

	rv := reflect.ValueOf(value)
	if ishelper.IsNil(rv) || rv.IsZero() {
		return &Violation{
			Code:    ViolationRequired,
			Message: formatMessage(ViolationRequired, nil),
		}
	}

	return nil
}
