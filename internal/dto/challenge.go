package dto

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

type Challenge struct {
	Challenge      []byte
	Difficulty     uint64
	SolutionLength uint64
}

func ChallengeUnmarshal(data []byte) (Challenge, error) {
	split := strings.SplitN(string(data), ":", 3)
	if len(split) != 3 {
		return Challenge{}, fmt.Errorf("invalid challenge message")
	}

	challenge, err := base64.StdEncoding.DecodeString(split[0])
	if err != nil {
		return Challenge{}, fmt.Errorf("decoding challenge: %w", err)
	}

	difficulty, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		return Challenge{}, fmt.Errorf("parsing difficulty: %w", err)
	}

	solutionLength, err := strconv.ParseInt(split[2], 10, 64)
	if err != nil {
		return Challenge{}, fmt.Errorf("parsing solution length: %w", err)
	}

	return Challenge{
		Challenge:      challenge,
		Difficulty:     uint64(difficulty),
		SolutionLength: uint64(solutionLength),
	}, nil
}

func ChallengeMarshal(c Challenge) []byte {
	return []byte(fmt.Sprintf(
		"%s:%d:%d",
		base64.StdEncoding.EncodeToString(c.Challenge),
		c.Difficulty,
		c.SolutionLength,
	))
}
