package rest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/devv4n/image-previewer/internal/config"
	"github.com/devv4n/image-previewer/internal/service"
)

type Server struct {
	Service *service.Service
	Config  *config.Config
	server  *http.Server
}

// NewServer создаёт новый server.
func NewServer(s *service.Service, cfg *config.Config) *Server {
	return &Server{Service: s, Config: cfg}
}

func (s *Server) Serve() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /fill/{width}/{height}/{img_url...}", s.PreviewHandler)

	s.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.Config.Port),
		Handler:           LogMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("server starting", "address", s.server.Addr)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("error starting server", "error", err)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("shutting down server gracefully")
	return s.server.Shutdown(ctx)
}
