package is

import (
	"context"
	"testing"
	"valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestMinLength(t *testing.T) {
	t.Parallel()

	rule := MinLength(2)
	require.Nil(t, rule(context.Background(), "ab"))
	require.Nil(t, rule(context.Background(), []int{1, 2}))
	require.Equal(t, ViolationMinLength, rule(context.Background(), "a").Code)
	require.Equal(t, ViolationMinLength, rule(context.Background(), []int{1}).Code)
	require.Equal(t, ViolationMinLength, rule(context.Background(), 123).Code)
	require.Nil(t, rule(context.Background(), ishelper.None[string]()))
	require.Nil(t, rule(context.Background(), ishelper.Some("abc")))
}
