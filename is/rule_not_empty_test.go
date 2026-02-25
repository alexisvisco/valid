package is

import (
	"context"
	"testing"
	"valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestNotEmpty(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	require.Nil(t, NotEmpty(ctx, "x"))
	require.Nil(t, NotEmpty(ctx, []int{1}))
	require.Equal(t, ViolationNotEmpty, NotEmpty(ctx, "").Code)
	require.Equal(t, ViolationNotEmpty, NotEmpty(ctx, []int{}).Code)
	require.Equal(t, ViolationNotEmpty, NotEmpty(ctx, 123).Code)
	require.Nil(t, NotEmpty(ctx, ishelper.None[string]()))
	require.Nil(t, NotEmpty(ctx, ishelper.Some("x")))
}
