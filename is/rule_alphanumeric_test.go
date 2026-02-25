package is

import (
	"context"
	"testing"
	"valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestAlphanumeric(t *testing.T) {
	t.Parallel()

	rule := Alphanumeric
	require.Nil(t, rule(context.Background(), "Abc123"))
	require.Equal(t, ViolationAlphaNum, rule(context.Background(), "abc_123").Code)
	require.Equal(t, ViolationAlphaNum, rule(context.Background(), 123).Code)
	require.Nil(t, rule(context.Background(), ishelper.None[string]()))
	require.Nil(t, rule(context.Background(), ishelper.Some("abc123")))
}
