package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jezzaho/goro-web/internal"
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
	s.config.readTimeout = 5 * time.Second
	s.config.writeTimeout = 5 * time.Second

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
		s.logger.Info("Server started on port ", s.config.port)
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

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if err := srv.Start(); err != nil {
		srv.logger.Error("Error during server startup: ", err)
		os.Exit(1)
	}

	srv.logger.Info("Server started on port: ", srv.config.port)

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
func (app *Application) MockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// From parsing
	err := r.ParseForm()
	if err != nil {
		log.Println("Error during Form Parsing: ", err)
	}
	log.Println(r.Form)
	// Access the query parameters
	// Carrier Format in two letters
	carrier := r.FormValue("carrier")
	dateFrom := r.FormValue("date-from")
	dateTo := r.FormValue("date-to")
	separate := r.FormValue("separate")

	// Carrier check for Query
	var carrierNumber int
	switch carrier {
	case "LH":
		carrierNumber = 0
	case "OS":
		carrierNumber = 1
	case "LX":
		carrierNumber = 2
	case "SN":
		carrierNumber = 3
	case "EN":
		carrierNumber = 4
	default:
		carrierNumber = 0
	}

	log.Println(carrierNumber)
	dateFromSSIM := internal.DateToSSIM(dateFrom)
	dateToSSIM := internal.DateToSSIM(dateTo)
	log.Println(dateFromSSIM)
	log.Println(dateToSSIM)
	separateBool := false
	if separate == "on" {
		separateBool = true
	}
	log.Println(separateBool)
	auth := internal.PostForAuth()
	query := internal.GetQueryListForAirline(carrierNumber, dateFromSSIM, dateToSSIM)
	data := internal.GetApiData(query, auth)
	data = internal.FlattenJSON(data)
	w.Header().Set("Content-Type", "text/csv")

	currentDate := time.Now().Format("20060102")
	filename := fmt.Sprintf("%s_%s.csv", currentDate, carrier)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Use the modified CreateCSVFromResponse function
	if err := internal.CreateCSVFromResponse(w, data, separateBool); err != nil {
		log.Printf("Error creating CSV: %v", err)
		http.Error(w, "Error creating CSV: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *Application) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "static/index.html")
		return
	}
	app.fs.ServeHTTP(w, r)
}
