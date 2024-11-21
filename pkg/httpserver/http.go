package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	defaultAdress          = "0.0.0.0:8080"
	defaultReadTimeout     = 3 * time.Second
	defaultWriteTimeout    = 3 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	server *http.Server

	shutdownTimeout time.Duration
}

func New(handler http.Handler, opts ...Option) *Server {
	httpserv := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAdress,
	}

	serv := &Server{
		server:          httpserv,
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(serv)
	}

	return serv
}

func (s *Server) Run(ctx context.Context, debugLog *logrus.Logger) error {
	debugLog.Debug(fmt.Sprintf("running http server on %v", s.server.Addr))

	go func() {
		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			debugLog.Debug("http server error:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	debugLog.Debug("Shutting down http server")

	ctx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
