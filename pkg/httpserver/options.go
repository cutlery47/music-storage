package httpserver

import (
	"fmt"
	"time"
)

type Option func(*Server)

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}

func Addr(host, port string) Option {
	return func(s *Server) {
		s.server.Addr = fmt.Sprintf("%v:%v", host, port)
	}
}
