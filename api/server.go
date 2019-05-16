package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/hasyimibhar/github-search/common"
	"github.com/hasyimibhar/github-search/github"
	"upper.io/db.v3/lib/sqlbuilder"
)

type Config struct {
	CORSEnabled bool
	HTTPPort    int
}

// Server handles API requests.
type Server struct {
	config       *Config
	log          common.Logger
	httpServer   *http.Server
	doneCh       chan struct{}
	shutdown     bool
	shutdownLock sync.Mutex
}

// NewServer creates an API server which immediately listens for requests.
func NewServer(cfg *Config, githubClient *github.Client, database sqlbuilder.Database, log common.Logger) (*Server, error) {

	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))

	if cfg.CORSEnabled {
		log.Info("api: CORS enabled")

		cors := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		})
		r.Use(cors.Handler)
	}

	router := Router(githubClient, database, log)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: router,
	}

	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		log.Infof("api: started server on [::]:%d", cfg.HTTPPort)

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(fmt.Errorf("api: %s", err))
		}
	}()

	return &Server{
		config:     cfg,
		log:        log,
		httpServer: httpServer,
		doneCh:     doneCh,
		shutdown:   false,
	}, nil
}

// Shutdown gracefully shuts down the API server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()

	if s.shutdown {
		return nil
	}

	s.httpServer.Shutdown(ctx)
	<-s.doneCh

	s.log.Info("api: shutdown complete")
	s.shutdown = true
	return nil
}
