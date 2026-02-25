package is

import (
	"context"
	"testing"
	"github.com/alexisvisco/valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestLessThan(t *testing.T) {
	t.Parallel()

	rule := LessThan(3)
	require.Nil(t, rule(context.Background(), 2))
	require.Equal(t, ViolationLT, rule(context.Background(), 3).Code)
	require.Equal(t, ViolationLT, rule(context.Background(), 4).Code)
	require.Equal(t, ViolationLT, rule(context.Background(), "3").Code)
	require.Nil(t, rule(context.Background(), ishelper.None[int]()))
	require.Nil(t, rule(context.Background(), ishelper.Some(2)))
}
