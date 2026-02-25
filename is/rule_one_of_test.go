package is

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"valid/ishelper"
)

func TestOneOf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		allowed   any // int or string slice, dispatched below
		value     any
		wantError bool
		wantCode  ViolationCode
	}{
		{name: "int match", allowed: []int{1, 2, 3}, value: 2, wantError: false},
		{name: "int no match", allowed: []int{1, 2, 3}, value: 4, wantError: true, wantCode: ViolationOneOf},
		{name: "string match", allowed: []string{"a", "b", "c"}, value: "b", wantError: false},
		{name: "string no match", allowed: []string{"a", "b", "c"}, value: "d", wantError: true, wantCode: ViolationOneOf},
		{name: "wrong type string", allowed: []int{1, 2, 3}, value: "x", wantError: true, wantCode: ViolationOneOf},
		{name: "wrong type struct", allowed: []int{1, 2, 3}, value: struct{ Tags []string }{}, wantError: true, wantCode: ViolationOneOf},
		{name: "nil value", allowed: []int{1, 2, 3}, value: nil, wantError: true, wantCode: ViolationOneOf},
		// Optional
		{name: "None[int]", allowed: []int{1, 2, 3}, value: ishelper.None[int](), wantError: false},
		{name: "Some match", allowed: []int{1, 2, 3}, value: ishelper.Some(2), wantError: false},
		{name: "Some no match", allowed: []int{1, 2, 3}, value: ishelper.Some(4), wantError: true, wantCode: ViolationOneOf},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got *Violation
			switch allowed := tt.allowed.(type) {
			case []int:
				got = OneOf(allowed...)(context.Background(), tt.value)
			case []string:
				got = OneOf(allowed...)(context.Background(), tt.value)
			default:
				require.Failf(t, "unsupported allowed type", "%T", tt.allowed)
			}

			require.Equal(t, tt.wantError, got != nil, "OneOf(%v)(%#v)", tt.allowed, tt.value)
			if tt.wantError {
				require.Equal(t, tt.wantCode, got.Code, "OneOf(%v)(%#v)", tt.allowed, tt.value)
			}
		})
	}
}
