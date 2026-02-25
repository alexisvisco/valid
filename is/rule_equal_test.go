package is

import (
	"context"
	"testing"
	"github.com/alexisvisco/valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestEqual(t *testing.T) {
	t.Parallel()

	t.Run("int", func(t *testing.T) {
		t.Parallel()
		rule := Equal(3)
		require.Nil(t, rule(context.Background(), 3))
		require.Equal(t, ViolationEQ, rule(context.Background(), 2).Code)
		require.Equal(t, ViolationEQ, rule(context.Background(), 4).Code)
		require.Equal(t, ViolationEQ, rule(context.Background(), "3").Code)
		require.Nil(t, rule(context.Background(), ishelper.None[int]()))
		require.Nil(t, rule(context.Background(), ishelper.Some(3)))
	})

	t.Run("string", func(t *testing.T) {
		t.Parallel()
		rule := Equal("ok")
		require.Nil(t, rule(context.Background(), "ok"))
		require.Equal(t, ViolationEQ, rule(context.Background(), "ko").Code)
		require.Equal(t, ViolationEQ, rule(context.Background(), 1).Code)
		require.Nil(t, rule(context.Background(), ishelper.None[string]()))
		require.Nil(t, rule(context.Background(), ishelper.Some("ok")))
	})

	t.Run("bool", func(t *testing.T) {
		t.Parallel()
		rule := Equal(true)
		require.Nil(t, rule(context.Background(), true))
		require.Equal(t, ViolationEQ, rule(context.Background(), false).Code)
	})
}
