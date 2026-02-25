package is

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"valid/ishelper"
)

func TestRequired(t *testing.T) {
	t.Parallel()

	var nilPtr *string
	var nilSlice []string
	var nilAny any

	tests := []struct {
		name      string
		value     any
		wantError bool
		wantCode  ViolationCode
	}{
		{name: "zero int", value: 0, wantError: true, wantCode: ViolationRequired},
		{name: "zero string", value: "", wantError: true, wantCode: ViolationRequired},
		{name: "zero bool", value: false, wantError: true, wantCode: ViolationRequired},
		{name: "nil untyped", value: nilAny, wantError: true, wantCode: ViolationRequired},
		{name: "nil pointer", value: nilPtr, wantError: true, wantCode: ViolationRequired},
		{name: "nil slice", value: nilSlice, wantError: true, wantCode: ViolationRequired},
		{name: "zero struct", value: struct{ A int }{}, wantError: true, wantCode: ViolationRequired},
		{name: "non-zero int", value: 1, wantError: false},
		{name: "non-zero string", value: "x", wantError: false},
		{name: "true bool", value: true, wantError: false},
		{name: "non-nil empty slice", value: []string{}, wantError: false},
		// Optional
		{name: "None[string]", value: ishelper.None[string](), wantError: true, wantCode: ViolationRequired},
		{name: "Some empty string", value: ishelper.Some(""), wantError: false},
		{name: "Some non-empty string", value: ishelper.Some("x"), wantError: false},
		{name: "Some zero int", value: ishelper.Some(0), wantError: false},
		{name: "Some non-zero int", value: ishelper.Some(1), wantError: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Required(context.Background(), tt.value)
			require.Equal(t, tt.wantError, got != nil, "Required(%#v)", tt.value)
			if tt.wantError {
				require.Equal(t, tt.wantCode, got.Code, "Required(%#v)", tt.value)
			}
		})
	}
}
