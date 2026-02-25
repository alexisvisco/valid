package is

import (
	"context"
	"testing"
	"valid/ishelper"

	"github.com/stretchr/testify/require"
)

func TestMatches(t *testing.T) {
	t.Parallel()

	rule := Matches(`^[A-Z]{2}[0-9]{3}$`)
	require.Nil(t, rule(context.Background(), "AB123"))
	require.Equal(t, ViolationMatches, rule(context.Background(), "ab123").Code)
	require.Equal(t, ViolationMatches, rule(context.Background(), 123).Code)
	require.Nil(t, rule(context.Background(), ishelper.None[string]()))
	require.Nil(t, rule(context.Background(), ishelper.Some("AB123")))
}

func TestMatchesInvalidPatternPanic(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		_ = Matches(`([`)
	})
}
