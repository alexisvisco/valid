package is

import (
	"context"
	"testing"
	"github.com/alexisvisco/valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestGreaterThan(t *testing.T) {
	t.Parallel()

	rule := GreaterThan(3)
	require.Equal(t, ViolationGT, rule(context.Background(), 3).Code)
	require.Equal(t, ViolationGT, rule(context.Background(), 2).Code)
	require.Nil(t, rule(context.Background(), 4))
	require.Equal(t, ViolationGT, rule(context.Background(), "3").Code)
	require.Nil(t, rule(context.Background(), ishelper.None[int]()))
	require.Nil(t, rule(context.Background(), ishelper.Some(4)))
}
