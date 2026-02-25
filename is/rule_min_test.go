package is

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"valid/ishelper"
)

func TestMin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		min       any
		value     any
		wantError bool
		wantCode  ViolationCode
	}{
		{name: "int below", min: 1, value: 0, wantError: true, wantCode: ViolationMin},
		{name: "int equal", min: 1, value: 1, wantError: false},
		{name: "int above", min: 1, value: 2, wantError: false},
		{name: "float below", min: 1.0, value: 0.5, wantError: true, wantCode: ViolationMin},
		{name: "float equal", min: 1.0, value: 1.0, wantError: false},
		{name: "float above", min: 1.0, value: 1.5, wantError: false},
		{name: "uint below", min: uint(1), value: uint(0), wantError: true, wantCode: ViolationMin},
		{name: "uint equal", min: uint(1), value: uint(1), wantError: false},
		{name: "uint above", min: uint(1), value: uint(2), wantError: false},
		{name: "non numeric", min: 1, value: "x", wantError: true, wantCode: ViolationMin},
		{name: "nil value", min: 1, value: nil, wantError: true, wantCode: ViolationMin},
		{name: "max uint64 equal", min: uint64(math.MaxUint64), value: uint64(math.MaxUint64), wantError: false},
		{name: "max uint64 minus one", min: uint64(math.MaxUint64), value: uint64(math.MaxUint64 - 1), wantError: true, wantCode: ViolationMin},
		// Optional
		{name: "None[int]", min: 1, value: ishelper.None[int](), wantError: false},
		{name: "Some below", min: 1, value: ishelper.Some(0), wantError: true, wantCode: ViolationMin},
		{name: "Some equal", min: 1, value: ishelper.Some(1), wantError: false},
		{name: "Some above", min: 1, value: ishelper.Some(5), wantError: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got *Violation
			switch min := tt.min.(type) {
			case int:
				got = Min(min)(context.Background(), tt.value)
			case float64:
				got = Min(min)(context.Background(), tt.value)
			case uint:
				got = Min(min)(context.Background(), tt.value)
			case uint64:
				got = Min(min)(context.Background(), tt.value)
			default:
				require.Failf(t, "unsupported min type", "%T", tt.min)
			}

			require.Equal(t, tt.wantError, got != nil, "Min(%v)(%#v)", tt.min, tt.value)
			if tt.wantError {
				require.Equal(t, tt.wantCode, got.Code, "Min(%v)(%#v)", tt.min, tt.value)
			}
		})
	}
}
