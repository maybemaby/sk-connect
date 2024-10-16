package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/maybemaby/sk-connect/gen/proto/api/v1/apiv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	mux          *http.ServeMux
	srv          *http2.Server
	logger       *slog.Logger
	Addr         string
	allowedHosts []string
	ctx          context.Context
}

type ServerConfig struct {
	Port         string
	LogLevel     slog.Level
	AllowedHosts []string
}

func NewServer(cfg ServerConfig) *Server {

	logger := BootstrapLogger(cfg.LogLevel, TextFormat, true)

	return &Server{
		mux:          http.NewServeMux(),
		srv:          &http2.Server{},
		logger:       logger,
		Addr:         ":" + cfg.Port,
		allowedHosts: cfg.AllowedHosts,
		ctx:          context.Background(),
	}
}

func (s *Server) MountHandlers() {
	rootMw := RootMiddleware(s.logger, s.allowedHosts[0])

	s.mux.Handle("/health", rootMw.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	sampleHandler := &SampleHandler{}
	samplePath, sampleRpc := apiv1connect.NewSampleServiceHandler(sampleHandler)
	s.logger.Debug("Mounting sample handler", slog.String("path", samplePath))
	s.mux.Handle(samplePath, rootMw.Then(sampleRpc))
}

func (s *Server) Start() error {
	s.MountHandlers()

	s.logger.Info("Starting server at", slog.String("addr", s.Addr))

	return http.ListenAndServe(s.Addr, h2c.NewHandler(s.mux, s.srv))
}
