package pow

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/bits"
)

const uint8Size = 8

func NewChallenge(challengeLength uint64) ([]byte, error) {
	challenge := make([]byte, challengeLength)
	if _, err := rand.Read(challenge); err != nil {
		return nil, fmt.Errorf("generating challenge: %w", err)
	}

	return challenge, nil
}

func SolveChallenge(challenge []byte, solutionLength, difficulty uint64) ([]byte, error) {
	solution := make([]byte, solutionLength)
	for i := uint64(0); i < math.MaxUint64; i++ {
		binary.BigEndian.PutUint64(solution, i)
		verified, err := VerifyChallenge(challenge, solution, difficulty)
		if err != nil {
			return nil, fmt.Errorf("verifying solution: %w", err)
		}

		if verified {
			return solution, nil
		}
	}

	return nil, nil
}

func VerifyChallenge(challenge, solution []byte, difficulty uint64) (bool, error) {
	hasher := sha256.New()
	_, err := hasher.Write(challenge)
	if err != nil {
		return false, fmt.Errorf("writing challenge random to hasher: %w", err)
	}

	_, err = hasher.Write(solution)
	if err != nil {
		return false, fmt.Errorf("writing solution to hasher: %w", err)
	}

	hash := hasher.Sum(nil)
	return verifyNumberOfLeadingZeros(hash, difficulty), nil
}

func verifyNumberOfLeadingZeros(data []byte, zeros uint64) bool {
	var count uint64
	for _, b := range data {
		leadingZeros := bits.LeadingZeros8(b)
		count += uint64(leadingZeros)
		if count >= zeros {
			return true
		}

		if leadingZeros != uint8Size {
			return false
		}
	}

	return false
}
