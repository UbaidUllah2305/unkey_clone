package zen

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/unkeyed/unkey/go/pkg/fault"
	"github.com/unkeyed/unkey/go/pkg/otel/logging"
	"github.com/unkeyed/unkey/go/pkg/tls"
)

// Server manages HTTP server configuration, route registration, and lifecycle.
// It provides connection pooling for session objects to reduce memory churn
// during request handling.
//
// Server instances should be created with the New function and can be safely
// used by multiple goroutines.
type Server struct {
	mu sync.Mutex

	logger      logging.Logger
	isListening bool
	mux         *http.ServeMux
	srv         *http.Server
	flags       Flags
	tlsConfig   *tls.Config

	sessions sync.Pool
}

// Flags configures the behavior of a Server instance.
type Flags struct {
	// TestMode enables test mode, accepting certain headers from untrusted clients such as fake times for testing purposes.
	TestMode bool
}

// Config configures the behavior of a Server instance.
type Config struct {
	// Logger provides structured logging for the server. If nil, logging is disabled.
	Logger logging.Logger

	// TLS configuration for HTTPS connections.
	// If this is provided, the server will use HTTPS.
	TLS *tls.Config

	Flags *Flags
}

// New creates a new server with the provided configuration.
// It initializes the HTTP server and session pool with default timeouts.
//
// The HTTP server is configured with reasonable defaults for production use:
// - ReadTimeout: 10 seconds
// - WriteTimeout: 20 seconds
//
// Example:
//
//	server, err := zen.New(zen.Config{
//	    InstanceID: "api-server-1",
//	    Logger: logger,
//	})
//	if err != nil {
//	    log.Fatalf("failed to initialize server: %v", err)
//	}
func New(config Config) (*Server, error) {
	mux := http.NewServeMux()

	srv := &http.Server{
		Handler: mux,
		// See https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		//
		// > # http.ListenAndServe is doing it wrong
		// > Incidentally, this means that the package-level convenience functions that bypass http.Server
		// > like http.ListenAndServe, http.ListenAndServeTLS and http.Serve are unfit for public Internet
		// > Servers.
		// >
		// > Those functions leave the Timeouts to their default off value, with no way of enabling them,
		// > so if you use them you'll soon be leaking connections and run out of file descriptors. I've
		// > made this mistake at least half a dozen times.
		// >
		// > Instead, create a http.Server instance with ReadTimeout and WriteTimeout and use its
		// > corresponding methods, like in the example a few paragraphs above.
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	flags := Flags{
		TestMode: false,
	}
	if config.Flags != nil {
		flags = *config.Flags
	}
	s := &Server{
		mu:          sync.Mutex{},
		logger:      config.Logger,
		isListening: false,
		mux:         mux,
		srv:         srv,
		flags:       flags,
		tlsConfig:   config.TLS,
		sessions: sync.Pool{
			New: func() any {
				return &Session{
					WorkspaceID:    "",
					requestID:      "",
					w:              nil,
					r:              nil,
					requestBody:    []byte{},
					responseStatus: 0,
					responseBody:   []byte{},
				}
			},
		},
	}
	return s, nil
}

// Get a fresh or reused session from the pool.
//
// You must return it to the pool when you are done via s.returnSession()
//
// Returning proper types is a bit annoying here, cause each session might have
// different requets and response types, so we just return any.
// You should immediately cast it to your desired type
//
// sess := s.getSession().(session.Session[MyRequest, MyResponse]).
func (s *Server) getSession() any {
	return s.sessions.Get()
}

// Return the session to the sync pool.
func (s *Server) returnSession(session any) {
	s.sessions.Put(session)
}

// Mux returns the underlying http.ServeMux.
// This is primarily intended for testing and advanced usage scenarios.
func (s *Server) Mux() *http.ServeMux {
	return s.mux
}

func (s *Server) Flags() Flags {
	return s.flags
}

// Listen starts the HTTP server on the specified address.
// This method blocks until the server shuts down or encounters an error.
// Once listening, the server will not start again if Listen is called multiple times.
// If TLS configuration is provided, the server will use HTTPS.
//
// The provided context is used to gracefully shut down the server when the context is canceled.
//
// Example:
//
//	// Start server in a goroutine to allow for graceful shutdown
//	go func() {
//	    if err := server.Listen(ctx, ":8080"); err != nil {
//	        log.Printf("server stopped: %v", err)
//	    }
//	}()
func (s *Server) Serve(ctx context.Context, ln net.Listener) error {
	s.mu.Lock()
	if s.isListening {
		s.logger.Warn("already listening")
		s.mu.Unlock()
		return nil
	}
	s.isListening = true
	s.mu.Unlock()

	// Set up context handling for graceful shutdown
	serverCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle context cancellation in a separate goroutine
	go func() {
		select {
		case <-ctx.Done():
			s.logger.Info("shutting down server due to context cancellation")
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer shutdownCancel()

			err := s.Shutdown(shutdownCtx)
			if err != nil {
				s.logger.Error("error during server shutdown", "error", err.Error())
			}
		case <-serverCtx.Done():
			// Server stopped on its own
		}
	}()
	var err error

	// Check if TLS should be used
	if s.tlsConfig != nil {
		s.logger.Info("listening", "srv", "https", "addr", ln.Addr().String())

		s.srv.TLSConfig = s.tlsConfig

		// ListenAndServeTLS with empty strings will use the certificates from TLSConfig
		err = s.srv.ServeTLS(ln, "", "")
	} else {
		s.logger.Info("listening", "srv", "http", "addr", ln.Addr().String())
		err = s.srv.Serve(ln)
	}

	// Cancel the server context since the server has stopped
	cancel()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fault.Wrap(err, fault.Internal("listening failed"))
	}
	return nil
}

// RegisterRoute adds an HTTP route to the server with the specified middleware chain.
// Routes are matched by both method and path.
//
// Middleware is applied in the order provided, with each middleware wrapping the next.
// The innermost handler (last to execute) is the route's handler.
//
// Example:
//
//	server.RegisterRoute(
//	    []zen.Middleware{zen.WithLogging(logger), zen.WithErrorHandling()},
//	    zen.NewRoute("GET", "/health", healthCheckHandler),
//	)
func (s *Server) RegisterRoute(middlewares []Middleware, route Route) {
	s.logger.Info("registering",
		"method", route.Method(),
		"path", route.Path(),
	)
	s.mux.HandleFunc(
		fmt.Sprintf("%s %s", route.Method(), route.Path()),
		func(w http.ResponseWriter, r *http.Request) {
			sess, ok := s.getSession().(*Session)
			if !ok {
				panic("Unable to cast session")
			}
			defer func() {
				sess.reset()
				s.returnSession(sess)
			}()

			err := sess.init(w, r)
			if err != nil {
				s.logger.Error("failed to init session")
				return
			}

			// Apply middleware
			var handle HandleFunc = route.Handle

			// Reverses the middlewares to run in the desired order.
			// If middlewares are [A, B, C], this writes [C, B, A] to s.middlewares.
			for i := len(middlewares) - 1; i >= 0; i-- {
				handle = middlewares[i](handle)
			}

			err = handle(r.Context(), sess)

			if err != nil {
				panic(err)
			}
		})
}

// Shutdown gracefully stops the HTTP server, allowing in-flight requests
// to complete before returning or the context is canceled.
//
// Example:
//
//	// Handle shutdown signal
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	if err := server.Shutdown(ctx); err != nil {
//	    log.Printf("server shutdown error: %v", err)
//	}
func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	s.isListening = false
	s.mu.Unlock()

	err := s.srv.Shutdown(ctx)
	if err != nil {
		return fault.Wrap(err)
	}
	return nil
}
