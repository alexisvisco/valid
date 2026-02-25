package is

import (
	"context"
	"testing"
	"valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestUUID(t *testing.T) {
	t.Parallel()

	rule := UUID
	require.Nil(t, rule(context.Background(), "123e4567-e89b-12d3-a456-426614174000"))
	require.Nil(t, rule(context.Background(), "123E4567-E89B-12D3-A456-426614174000"))
	require.Equal(t, ViolationUUID, rule(context.Background(), "bad-uuid").Code)
	require.Nil(t, rule(context.Background(), ishelper.None[string]()))
	require.Nil(t, rule(context.Background(), ishelper.Some("123e4567-e89b-12d3-a456-426614174000")))
}
