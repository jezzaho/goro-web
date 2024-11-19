package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

type Application struct {
	fs http.Handler
}

type AppLogger struct{}

func (l *AppLogger) Info(msg string, fields ...interface{})  { log.Printf("INFO: "+msg, fields...) }
func (l *AppLogger) Error(msg string, fields ...interface{}) { log.Printf("ERROR: "+msg, fields...) }

type Server struct {
	server *http.Server
	router *http.ServeMux

	config struct {
		port         string
		readTimeout  time.Duration
		writeTimeout time.Duration
	}

	// State management
	wg      sync.WaitGroup
	closing chan struct{}

	errors chan error

	logger AppLogger
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		router:  http.NewServeMux(),
		closing: make(chan struct{}),
		errors:  make(chan error),
	}
	s.config.port = ":3333"
	s.config.readTimeout = 35 * time.Second
	s.config.writeTimeout = 50 * time.Second

	for _, opt := range opts {
		opt(s)
	}

	s.server = &http.Server{
		Addr:         s.config.port,
		Handler:      s.router,
		ReadTimeout:  s.config.readTimeout,
		WriteTimeout: s.config.writeTimeout,
	}
	s.logger = AppLogger{}

	return s
}

func (s *Server) Start() error {
	var handler http.Handler = s.router
	s.server.Handler = handler

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.errors <- err
		}

	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	close(s.closing)

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	s.wg.Wait()
	return nil
}

type Option func(*Server)

func WithPort(port string) Option {
	return func(s *Server) {
		s.config.port = port
	}
}

func main() {
	godotenv.Load()

	srv := NewServer(WithPort(":3333"))

	app := Application{}

	fs := http.FileServer(http.Dir("static"))
	app.fs = fs

	srv.router.HandleFunc("/", app.IndexHandler)
	srv.router.HandleFunc("/csv", app.MockHandler)
	srv.router.HandleFunc("/progress", app.ProgressStreamHandler)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if err := srv.Start(); err != nil {
		srv.logger.Error("Error during server startup: ", err)
		os.Exit(1)
	}
	srv.logger.Info("Server started on port ", srv.config.port)

	sig := <-signalChan
	srv.logger.Info("Received shutdown signal: ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		srv.logger.Error("Graceful Server shutdown failed", err)
		os.Exit(1)
	}

	srv.logger.Info("Server shutdown completed")

}
