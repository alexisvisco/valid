package is

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"valid/ishelper"
)

func TestMax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		max       any
		value     any
		wantError bool
		wantCode  ViolationCode
	}{
		{name: "int below", max: 5, value: 4, wantError: false},
		{name: "int equal", max: 5, value: 5, wantError: false},
		{name: "int above", max: 5, value: 6, wantError: true, wantCode: ViolationMax},
		{name: "float below", max: 5.0, value: 4.5, wantError: false},
		{name: "float equal", max: 5.0, value: 5.0, wantError: false},
		{name: "float above", max: 5.0, value: 5.5, wantError: true, wantCode: ViolationMax},
		{name: "uint below", max: uint(5), value: uint(4), wantError: false},
		{name: "uint equal", max: uint(5), value: uint(5), wantError: false},
		{name: "uint above", max: uint(5), value: uint(6), wantError: true, wantCode: ViolationMax},
		{name: "non numeric", max: 5, value: "x", wantError: true, wantCode: ViolationMax},
		{name: "nil value", max: 5, value: nil, wantError: true, wantCode: ViolationMax},
		{name: "max uint64 equal", max: uint64(math.MaxUint64), value: uint64(math.MaxUint64), wantError: false},
		// Optional
		{name: "None[int]", max: 5, value: ishelper.None[int](), wantError: false},
		{name: "Some below", max: 5, value: ishelper.Some(3), wantError: false},
		{name: "Some equal", max: 5, value: ishelper.Some(5), wantError: false},
		{name: "Some above", max: 5, value: ishelper.Some(6), wantError: true, wantCode: ViolationMax},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got *Violation
			switch max := tt.max.(type) {
			case int:
				got = Max(max)(context.Background(), tt.value)
			case float64:
				got = Max(max)(context.Background(), tt.value)
			case uint:
				got = Max(max)(context.Background(), tt.value)
			case uint64:
				got = Max(max)(context.Background(), tt.value)
			default:
				require.Failf(t, "unsupported max type", "%T", tt.max)
			}

			require.Equal(t, tt.wantError, got != nil, "Max(%v)(%#v)", tt.max, tt.value)
			if tt.wantError {
				require.Equal(t, tt.wantCode, got.Code, "Max(%v)(%#v)", tt.max, tt.value)
			}
		})
	}
}
