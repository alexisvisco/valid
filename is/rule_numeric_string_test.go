package is

import (
	"context"
	"testing"
	"valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestNumeric(t *testing.T) {
	t.Parallel()

	rule := Numeric
	require.Nil(t, rule(context.Background(), "123"))
	require.Nil(t, rule(context.Background(), "-123.45"))
	require.Equal(t, ViolationNumeric, rule(context.Background(), "12a3").Code)
	require.Equal(t, ViolationNumeric, rule(context.Background(), 123).Code)
	require.Nil(t, rule(context.Background(), ishelper.None[string]()))
	require.Nil(t, rule(context.Background(), ishelper.Some("12")))
}
