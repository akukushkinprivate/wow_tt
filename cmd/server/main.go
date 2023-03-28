package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/akukushkinprivate/wow_tt/internal/handler"
	"github.com/akukushkinprivate/wow_tt/internal/server"
	"github.com/caarlos0/env/v7"
)

func main() {
	cfg := mustLoadConfig()
	wowHandler := handler.New(cfg.ChallengeLength, cfg.SolutionLength, cfg.ChallengeDifficulty, cfg.ConnTimeout)
	srv := server.MustNewServer(cfg.Addr, wowHandler)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		srv.ListenAndServe()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	log.Println("server listening", cfg.Addr)
	<-c
	log.Println("server shutting down...")
	srv.Stop()
	wg.Wait()
}

func mustLoadConfig() config {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("failed to parse env config:", err)
	}

	return cfg
}

type config struct {
	Addr                string        `env:"ADDR" envDefault:":80"`
	ChallengeLength     uint64        `env:"CHALLENGE_LENGTH" envDefault:"32"`
	SolutionLength      uint64        `env:"SOLUTION_LENGTH" envDefault:"8"`
	ChallengeDifficulty uint64        `env:"CHALLENGE_DIFFICULTY" envDefault:"20"`
	ConnTimeout         time.Duration `env:"CONN_TIMEOUT" envDefault:"1m"`
}
