package is

import (
	"context"
	"testing"
	"github.com/alexisvisco/valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestEmail(t *testing.T) {
	t.Parallel()

	require.Nil(t, Email(context.Background(), "user@example.com"))
	require.Equal(t, ViolationEmail, Email(context.Background(), "not-an-email").Code)
	require.Equal(t, ViolationEmail, Email(context.Background(), 12).Code)
	require.Nil(t, Email(context.Background(), ishelper.None[string]()))
	require.Nil(t, Email(context.Background(), ishelper.Some("user@example.com")))
}
