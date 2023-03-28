package pow

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSolveChallenge(t *testing.T) {
	challenge, err := NewChallenge(32)
	require.NoError(t, err)

	solution, err := SolveChallenge(challenge, 8, 20)
	require.NoError(t, err)

	verified, err := VerifyChallenge(challenge, solution, 20)
	require.NoError(t, err)

	require.True(t, verified)
}
