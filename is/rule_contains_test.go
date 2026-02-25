package is

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/alexisvisco/valid/ishelper"
)

func TestContains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		elem      any // string or int, dispatched below
		value     any
		wantError bool
		wantCode  ViolationCode
	}{
		// String contains substring
		{name: "string contains", elem: "foo", value: "foobar", wantError: false},
		{name: "string contains exact", elem: "foo", value: "foo", wantError: false},
		{name: "string not contains", elem: "foo", value: "barbaz", wantError: true, wantCode: ViolationContains},
		{name: "string empty elem", elem: "", value: "anything", wantError: false},
		// Slice contains element
		{name: "int slice contains", elem: 3, value: []int{1, 2, 3}, wantError: false},
		{name: "int slice not contains", elem: 4, value: []int{1, 2, 3}, wantError: true, wantCode: ViolationContains},
		{name: "string slice contains", elem: "b", value: []string{"a", "b", "c"}, wantError: false},
		{name: "string slice not contains", elem: "d", value: []string{"a", "b", "c"}, wantError: true, wantCode: ViolationContains},
		{name: "empty slice", elem: 1, value: []int{}, wantError: true, wantCode: ViolationContains},
		// Unsupported
		{name: "non-string non-slice value", elem: "foo", value: 42, wantError: true, wantCode: ViolationContains},
		{name: "nil value", elem: "foo", value: nil, wantError: true, wantCode: ViolationContains},
		// Optional
		{name: "None[string]", elem: "foo", value: ishelper.None[string](), wantError: false},
		{name: "Some string contains", elem: "foo", value: ishelper.Some("foobar"), wantError: false},
		{name: "Some string not contains", elem: "foo", value: ishelper.Some("bar"), wantError: true, wantCode: ViolationContains},
		{name: "None[[]int]", elem: 1, value: ishelper.None[[]int](), wantError: false},
		{name: "Some slice contains", elem: 2, value: ishelper.Some([]int{1, 2, 3}), wantError: false},
		{name: "Some slice not contains", elem: 4, value: ishelper.Some([]int{1, 2, 3}), wantError: true, wantCode: ViolationContains},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got *Violation
			switch elem := tt.elem.(type) {
			case string:
				got = Contains(elem)(context.Background(), tt.value)
			case int:
				got = Contains(elem)(context.Background(), tt.value)
			default:
				require.Failf(t, "unsupported elem type", "%T", tt.elem)
			}

			require.Equal(t, tt.wantError, got != nil, "Contains(%v)(%#v)", tt.elem, tt.value)
			if tt.wantError {
				require.Equal(t, tt.wantCode, got.Code, "Contains(%v)(%#v)", tt.elem, tt.value)
			}
		})
	}
}
