package is

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"valid/ishelper"
)

func TestHasPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		prefix    string
		value     any
		wantError bool
		wantCode  ViolationCode
	}{
		{name: "match", prefix: "foo", value: "foobar", wantError: false},
		{name: "exact match", prefix: "foo", value: "foo", wantError: false},
		{name: "no match", prefix: "foo", value: "barfoo", wantError: true, wantCode: ViolationHasPrefix},
		{name: "empty string", prefix: "foo", value: "", wantError: true, wantCode: ViolationHasPrefix},
		{name: "empty prefix", prefix: "", value: "anything", wantError: false},
		{name: "non-string value", prefix: "foo", value: 42, wantError: true, wantCode: ViolationHasPrefix},
		{name: "nil value", prefix: "foo", value: nil, wantError: true, wantCode: ViolationHasPrefix},
		// Optional
		{name: "None[string]", prefix: "foo", value: ishelper.None[string](), wantError: false},
		{name: "Some match", prefix: "foo", value: ishelper.Some("foobar"), wantError: false},
		{name: "Some no match", prefix: "foo", value: ishelper.Some("barfoo"), wantError: true, wantCode: ViolationHasPrefix},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := HasPrefix(tt.prefix)(context.Background(), tt.value)
			require.Equal(t, tt.wantError, got != nil, "HasPrefix(%q)(%#v)", tt.prefix, tt.value)
			if tt.wantError {
				require.Equal(t, tt.wantCode, got.Code, "HasPrefix(%q)(%#v)", tt.prefix, tt.value)
			}
		})
	}
}
