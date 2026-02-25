package is

import (
	"context"
	"testing"
	"valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestAlpha(t *testing.T) {
	t.Parallel()

	rule := Alpha
	require.Nil(t, rule(context.Background(), "AbCd"))
	require.Equal(t, ViolationAlpha, rule(context.Background(), "abc1").Code)
	require.Equal(t, ViolationAlpha, rule(context.Background(), 123).Code)
	require.Nil(t, rule(context.Background(), ishelper.None[string]()))
	require.Nil(t, rule(context.Background(), ishelper.Some("abc")))
}
