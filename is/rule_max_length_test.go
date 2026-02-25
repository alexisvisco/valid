package is

import (
	"context"
	"testing"
	"valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestMaxLength(t *testing.T) {
	t.Parallel()

	rule := MaxLength(2)
	require.Nil(t, rule(context.Background(), "ab"))
	require.Nil(t, rule(context.Background(), []int{1, 2}))
	require.Equal(t, ViolationMaxLength, rule(context.Background(), "abc").Code)
	require.Equal(t, ViolationMaxLength, rule(context.Background(), []int{1, 2, 3}).Code)
	require.Equal(t, ViolationMaxLength, rule(context.Background(), 123).Code)
	require.Nil(t, rule(context.Background(), ishelper.None[string]()))
	require.Nil(t, rule(context.Background(), ishelper.Some("ab")))
}
