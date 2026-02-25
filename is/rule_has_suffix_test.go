package is

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"valid/ishelper"
)

func TestHasSuffix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		suffix    string
		value     any
		wantError bool
		wantCode  ViolationCode
	}{
		{name: "match", suffix: "bar", value: "foobar", wantError: false},
		{name: "exact match", suffix: "bar", value: "bar", wantError: false},
		{name: "no match", suffix: "bar", value: "barfoo", wantError: true, wantCode: ViolationHasSuffix},
		{name: "empty string", suffix: "bar", value: "", wantError: true, wantCode: ViolationHasSuffix},
		{name: "empty suffix", suffix: "", value: "anything", wantError: false},
		{name: "non-string value", suffix: "bar", value: 42, wantError: true, wantCode: ViolationHasSuffix},
		{name: "nil value", suffix: "bar", value: nil, wantError: true, wantCode: ViolationHasSuffix},
		// Optional
		{name: "None[string]", suffix: "bar", value: ishelper.None[string](), wantError: false},
		{name: "Some match", suffix: "bar", value: ishelper.Some("foobar"), wantError: false},
		{name: "Some no match", suffix: "bar", value: ishelper.Some("barfoo"), wantError: true, wantCode: ViolationHasSuffix},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := HasSuffix(tt.suffix)(context.Background(), tt.value)
			require.Equal(t, tt.wantError, got != nil, "HasSuffix(%q)(%#v)", tt.suffix, tt.value)
			if tt.wantError {
				require.Equal(t, tt.wantCode, got.Code, "HasSuffix(%q)(%#v)", tt.suffix, tt.value)
			}
		})
	}
}
