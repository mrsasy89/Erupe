package signv2server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	config "erupe-ce/config"
)

type Config struct {
	Logger      *zap.Logger
	DB          *sqlx.DB
	ErupeConfig *config.Config
}

// Server is the MHF custom launcher sign server.
type Server struct {
	sync.Mutex
	logger         *zap.Logger
	erupeConfig    *config.Config
	db             *sqlx.DB
	httpServer     *http.Server
	isShuttingDown bool
}

// NewServer creates a new Server type.
func NewServer(config *Config) *Server {
	s := &Server{
		logger:      config.Logger,
		erupeConfig: config.ErupeConfig,
		db:          config.DB,
		httpServer:  &http.Server{},
	}
	return s
}

// Start starts the server in a new goroutine.
func (s *Server) Start() error {
	// Set up the routes responsible for serving the launcher info and auth.
	r := mux.NewRouter()

	// Alias root to /launcher so GET / doesn’t 404
	r.HandleFunc("/", s.Launcher).Methods("GET")
	// Main launcher endpoint
	r.HandleFunc("/launcher", s.Launcher).Methods("GET")

	// Authentication and character endpoints
	r.HandleFunc("/login", s.Login)
	r.HandleFunc("/register", s.Register)
	r.HandleFunc("/character/create", s.CreateCharacter)
	r.HandleFunc("/character/delete", s.DeleteCharacter)
	r.HandleFunc("/character/export", s.ExportSave)

	// Wrap with CORS and logging
	handler := handlers.CORS(handlers.AllowedHeaders([]string{"Content-Type"}))(r)
	s.httpServer.Handler = handlers.LoggingHandler(os.Stdout, handler)

	// Listen on configured SignV2 port
	s.httpServer.Addr = fmt.Sprintf(":%d", s.erupeConfig.SignV2.Port)

	// Run server
	serveError := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			serveError <- err
		}
	}()

	// Wait a moment to confirm startup
	select {
	case err := <-serveError:
		return err
	case <-time.After(250 * time.Millisecond):
		return nil
	}
}

// Shutdown exits the server gracefully.
func (s *Server) Shutdown() {
	s.logger.Debug("Shutting down")

	s.Lock()
	s.isShuttingDown = true
	s.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Warn("Got error on httpServer shutdown", zap.Error(err))
	}
}