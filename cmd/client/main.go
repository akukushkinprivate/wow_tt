package main

import (
	"log"
	"net"
	"time"

	"github.com/akukushkinprivate/wow_tt/internal/dto"
	"github.com/akukushkinprivate/wow_tt/internal/pow"
	"github.com/caarlos0/env/v7"
)

const bufferCap = 1024

func main() {
	cfg := mustLoadConfig()

	tcpAddr, err := net.ResolveTCPAddr("tcp", cfg.Addr)
	if err != nil {
		log.Fatalln("failed to resolve server address:", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalln("failed to dial:", err)
	}

	defer func(conn net.Conn) {
		if err := conn.Close(); err != nil {
			log.Fatalln("failed to close connection to server:", err)
		}
	}(conn)

	buffer := make([]byte, bufferCap)
	read, err := conn.Read(buffer)
	if err != nil {
		log.Fatalln("failed to read challenge:", err)
	}

	log.Println("got challenge from server:", string(buffer[:read]))
	challenge, err := dto.ChallengeUnmarshal(buffer[:read])
	if err != nil {
		log.Fatalln("failed to unmarshal challenge:", err)
	}

	solution, err := pow.SolveChallenge(challenge.Challenge, challenge.SolutionLength, challenge.Difficulty)
	if err != nil {
		log.Fatalln("failed to solve challenge:", err)
	}

	if _, err := conn.Write(solution); err != nil {
		log.Fatalln("failed to write solution:", err)
	}

	log.Println("sent solution to server")

	read, err = conn.Read(buffer)
	if err != nil {
		log.Fatalln("failed to read quote:", err)
	}

	log.Println("got quote from server:", string(buffer[:read]))
}

func mustLoadConfig() config {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("failed to parse env config:", err)
	}

	return cfg
}

type config struct {
	Addr        string        `env:"SERVER_ADDR" envDefault:":80"`
	ConnTimeout time.Duration `env:"CONN_TIMEOUT" envDefault:"1m"`
}
