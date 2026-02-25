package is

import (
	"context"
	"testing"
	"valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestPositive(t *testing.T) {
	t.Parallel()

	rule := Positive
	require.Nil(t, rule(context.Background(), 1))
	require.Nil(t, rule(context.Background(), 0.001))
	require.Equal(t, ViolationPositive, rule(context.Background(), 0).Code)
	require.Equal(t, ViolationPositive, rule(context.Background(), -1).Code)
	require.Equal(t, ViolationPositive, rule(context.Background(), "1").Code)
	require.Nil(t, rule(context.Background(), ishelper.None[int]()))
	require.Nil(t, rule(context.Background(), ishelper.Some(2)))
}
