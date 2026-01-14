package web

import (
	"context"
	"expvar"
	"fmt"
	"habana_bmc_exporter/pkg/web/checkgrp"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

const DefaultPort = "4001"

type ServerOpts struct {
	Log        *logrus.Entry
	Port       string
	Middleware []Middleware
}

type Server struct {
	ServerOpts
	*http.Server
	mux *http.ServeMux
}

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request)

func NewServer(ctx context.Context, serverOpts ServerOpts, cgh checkgrp.Handlers, handlers []WebHandler, opts ...func(*Server)) *Server {
	if serverOpts.Port == "" {
		serverOpts.Port = DefaultPort
	}

	// Register readiness and liveness routes.
	mux := debugStandardLibraryMux()
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	server := &Server{
		ServerOpts: serverOpts,
		mux:        mux,
		Server: &http.Server{
			Addr:         fmt.Sprintf(":%s", serverOpts.Port),
			Handler:      mux,
			IdleTimeout:  1 * time.Minute,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	}

	// Register the user provided routes
	for _, h := range handlers {
		server.HandleFunc(ctx, h.Route, h.Handler, serverOpts.Log, serverOpts.Middleware...)
	}

	// override default config with user's provided params.
	for _, f := range opts {
		f(server)
	}
	return server
}

// HandleFunc will override the method HandleFunc of the http.Server.
func (s *Server) HandleFunc(ctx context.Context, path string, handler Handler, log *logrus.Entry, mw ...Middleware) {

	// Add the middleware to the handler.
	handler = wrapMiddleware(mw, handler, log)

	h := func(w http.ResponseWriter, r *http.Request) {

		// call the handler we passed.
		handler(ctx, w, r)
	}

	s.mux.HandleFunc(path, h)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

type WebHandler struct {
	Handler Handler
	Route   string
}

func WithReadTimeout(timeout time.Duration) func(*Server) {
	return func(s *Server) {
		s.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) func(*Server) {
	return func(s *Server) {
		s.WriteTimeout = timeout
	}
}

// debugStandardLibraryMux registers all the debug routes from the standard library
// into a new mux bypassing the use of the DefaultServerMux.
// Using the DefaultServerMux would be a security risk since a dependency could inject a handler
// into our service without us knowing it.
func debugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	// our metrics
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

type Closer interface {
	Close() error
}

type TlsCerts struct {
	CertFile string
	KeyFile  string
}

// Start will start the server and catch errors and signals.
func (s *Server) Start(log *logrus.Entry, cancelFunc context.CancelFunc, onClose Closer, certs ...TlsCerts) error {
	// Close the context so all request will be closed.

	defer cancelFunc()

	// Make a channel to listen for errors coming from the listener.
	// Use a buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	go func() {

		if len(certs) > 0 {
			if err := s.ListenAndServeTLS(certs[0].CertFile, certs[0].KeyFile); err != nil {
				serverErrors <- err
			}
		} else {
			if err := s.ListenAndServe(); err != nil {
				serverErrors <- err
			}
		}

	}()
	// buffered channel of 1 means sends happens before receive.
	shutdown := make(chan os.Signal, 1)
	// SIGINT - ctl +c
	// SIGTERM - k8s shutdown signal
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// ============================================================
	// Shutdown

	// Blocking main and wait for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.WithFields(logrus.Fields{
			"status": "shutdown started",
			"signal": sig,
		}).Info("shutdown")
		defer log.Info("shutdown complete")
		// Give outstanding requests a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if onClose != nil {
			err := onClose.Close()
			if err != nil {
				log.WithError(fmt.Errorf("failed closing: %w", err)).Error()
			}
		}

		cancelFunc()

		// Asking listener to shutdown and shed load.
		if err := s.Shutdown(ctx); err != nil {
			s.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
