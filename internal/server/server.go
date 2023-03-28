package server

import (
	"log"
	"net"
	"sync"
)

type handler interface {
	Handle(conn net.Conn)
}

type Server struct {
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
	handler  handler
}

func MustNewServer(addr string, handler handler) *Server {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	return &Server{
		listener: listener,
		quit:     make(chan interface{}),
		wg:       sync.WaitGroup{},
		handler:  handler,
	}
}

func (s *Server) ListenAndServe() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Println("failed to accept:", err)
			}

			return
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()

			s.handler.Handle(conn)
		}()
	}
}

func (s *Server) Stop() {
	close(s.quit)
	if err := s.listener.Close(); err != nil {
		log.Println("failed to close listener:", err)
	}

	s.wg.Wait()
}
