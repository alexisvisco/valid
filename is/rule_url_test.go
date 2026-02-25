package is

import (
	"context"
	"testing"
	"github.com/alexisvisco/valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestURL(t *testing.T) {
	t.Parallel()

	require.Nil(t, URL(context.Background(), "https://example.com/path"))
	require.Equal(t, ViolationURL, URL(context.Background(), "example.com/path").Code)
	require.Equal(t, ViolationURL, URL(context.Background(), 12).Code)
	require.Nil(t, URL(context.Background(), ishelper.None[string]()))
	require.Nil(t, URL(context.Background(), ishelper.Some("https://example.com")))
}
