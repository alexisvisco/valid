package is

import (
	"context"
	"testing"
	"github.com/alexisvisco/valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestNonNegative(t *testing.T) {
	t.Parallel()

	rule := NonNegative
	require.Nil(t, rule(context.Background(), 0))
	require.Nil(t, rule(context.Background(), 5))
	require.Nil(t, rule(context.Background(), 0.25))
	require.Equal(t, ViolationNonNeg, rule(context.Background(), -1).Code)
	require.Equal(t, ViolationNonNeg, rule(context.Background(), "-1").Code)
	require.Nil(t, rule(context.Background(), ishelper.None[int]()))
	require.Nil(t, rule(context.Background(), ishelper.Some(0)))
}
