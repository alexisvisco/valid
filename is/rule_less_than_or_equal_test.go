package is

import (
	"context"
	"testing"
	"github.com/alexisvisco/valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestLessThanOrEqual(t *testing.T) {
	t.Parallel()

	rule := LessThanOrEqual(3)
	require.Nil(t, rule(context.Background(), 2))
	require.Nil(t, rule(context.Background(), 3))
	require.Equal(t, ViolationLTE, rule(context.Background(), 4).Code)
	require.Equal(t, ViolationLTE, rule(context.Background(), "3").Code)
	require.Nil(t, rule(context.Background(), ishelper.None[int]()))
	require.Nil(t, rule(context.Background(), ishelper.Some(3)))
}
