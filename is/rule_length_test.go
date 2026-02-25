package is

import (
	"context"
	"testing"

	"valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		min       int
		max       int
		value     any
		wantError bool
		wantCode  ViolationCode
	}{
		{name: "string min", min: 1, max: 3, value: "a", wantError: false},
		{name: "string max", min: 1, max: 3, value: "abc", wantError: false},
		{name: "string below", min: 1, max: 3, value: "", wantError: true, wantCode: ViolationLength},
		{name: "string above", min: 1, max: 3, value: "abcd", wantError: true, wantCode: ViolationLength},
		{name: "slice in range", min: 1, max: 3, value: []int{1, 2}, wantError: false},
		{name: "slice below", min: 1, max: 3, value: []int{}, wantError: true, wantCode: ViolationLength},
		{name: "slice above", min: 1, max: 3, value: []int{1, 2, 3, 4}, wantError: true, wantCode: ViolationLength},
		{name: "map in range", min: 1, max: 3, value: map[string]int{"a": 1, "b": 2}, wantError: false},
		{name: "map below", min: 1, max: 3, value: map[string]int{}, wantError: true, wantCode: ViolationLength},
		{name: "map above", min: 1, max: 3, value: map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, wantError: true, wantCode: ViolationLength},
		{name: "unsupported type", min: 1, max: 3, value: 42, wantError: true, wantCode: ViolationLength},
		{name: "nil type", min: 1, max: 3, value: nil, wantError: true, wantCode: ViolationLength},
		// Optional
		{name: "None[string]", min: 1, max: 5, value: ishelper.None[string](), wantError: false},
		{name: "Some empty string", min: 1, max: 5, value: ishelper.Some(""), wantError: true, wantCode: ViolationLength},
		{name: "Some valid string", min: 1, max: 5, value: ishelper.Some("hi"), wantError: false},
		{name: "Some too long string", min: 1, max: 5, value: ishelper.Some("toolong"), wantError: true, wantCode: ViolationLength},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Length(tt.min, tt.max)(context.Background(), tt.value)
			require.Equal(t, tt.wantError, got != nil, "Length(%d,%d)(%#v)", tt.min, tt.max, tt.value)
			if tt.wantError {
				require.Equal(t, tt.wantCode, got.Code, "Length(%d,%d)(%#v)", tt.min, tt.max, tt.value)
			}
		})
	}
}
