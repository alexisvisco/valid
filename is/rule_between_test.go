package is

import (
	"context"
	"testing"
	"github.com/alexisvisco/valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestBetween(t *testing.T) {
	t.Parallel()

	rule := Between(10, 20)
	require.Nil(t, rule(context.Background(), 10))
	require.Nil(t, rule(context.Background(), 15))
	require.Nil(t, rule(context.Background(), 20))
	require.Equal(t, ViolationBetween, rule(context.Background(), 9).Code)
	require.Equal(t, ViolationBetween, rule(context.Background(), 21).Code)
	require.Equal(t, ViolationBetween, rule(context.Background(), "10").Code)
	require.Nil(t, rule(context.Background(), ishelper.None[int]()))
	require.Nil(t, rule(context.Background(), ishelper.Some(12)))
}
