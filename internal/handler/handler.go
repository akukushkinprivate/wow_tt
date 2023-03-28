package handler

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/akukushkinprivate/wow_tt/internal/dto"
	"github.com/akukushkinprivate/wow_tt/internal/pow"
	"github.com/akukushkinprivate/wow_tt/internal/quote"
)

type Handler struct {
	challengeLength     uint64
	solutionLength      uint64
	challengeDifficulty uint64
	connTimeout         time.Duration
}

func New(challengeLength, solutionLength, challengeDifficulty uint64, connTimeout time.Duration) *Handler {
	return &Handler{
		challengeLength:     challengeLength,
		solutionLength:      solutionLength,
		challengeDifficulty: challengeDifficulty,
		connTimeout:         connTimeout,
	}
}

func (h *Handler) Handle(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("failed to close connection:", err)
		}
	}(conn)

	ctx, cancel := context.WithTimeout(context.Background(), h.connTimeout)
	defer cancel()

	challenge, err := pow.NewChallenge(h.challengeLength)
	if err != nil {
		log.Println("failed to create new challenge:", err)
		return
	}

	challengeDTO := dto.Challenge{
		Challenge:      challenge,
		Difficulty:     h.challengeDifficulty,
		SolutionLength: h.solutionLength,
	}
	marshaled := dto.ChallengeMarshal(challengeDTO)
	if _, err = conn.Write(marshaled); err != nil {
		log.Println("failed to write challenge:", err)
		return
	}

	resultChannel := make(chan solutionChannel)
	go h.readSolution(conn, resultChannel)

	select {
	case <-ctx.Done():
		log.Println("failed to handle by timeout")
		return
	case result := <-resultChannel:
		if result.Err != nil {
			log.Println("failed to read solution:", err)
			return
		}

		verified, err := pow.VerifyChallenge(challenge, result.Solution, h.challengeDifficulty)
		if err != nil {
			log.Println("failed to verify solution:", err)
			return
		}

		if !verified {
			return
		}

		randomQuote := quote.GetRandomQuote()
		if _, err := conn.Write([]byte(randomQuote)); err != nil {
			log.Println("failed to write wow quote:", err)
		}
	}
}

func (h *Handler) readSolution(conn net.Conn, channel chan solutionChannel) {
	defer close(channel)

	solution := make([]byte, h.solutionLength)
	_, err := conn.Read(solution)
	if err != nil {
		channel <- solutionChannel{Err: err}
		return
	}

	channel <- solutionChannel{Solution: solution}
}

type solutionChannel struct {
	Solution []byte
	Err      error
}
