package is

import (
	"context"
	"testing"
	"github.com/alexisvisco/valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestGreaterThanOrEqual(t *testing.T) {
	t.Parallel()

	rule := GreaterThanOrEqual(3)
	require.Equal(t, ViolationGTE, rule(context.Background(), 2).Code)
	require.Nil(t, rule(context.Background(), 3))
	require.Nil(t, rule(context.Background(), 4))
	require.Equal(t, ViolationGTE, rule(context.Background(), "3").Code)
	require.Nil(t, rule(context.Background(), ishelper.None[int]()))
	require.Nil(t, rule(context.Background(), ishelper.Some(3)))
}
