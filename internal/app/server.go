package app

import (
	"github.com/efimovad/avito-internship/internal/store"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
	"net/http"
)

type Server struct {
	Mux          *mux.Router
	SessionStore sessions.Store
	Config       *Config
	Logger       *zap.SugaredLogger
	Sanitizer    *bluemonday.Policy
}

func NewServer(config *Config, logger *zap.SugaredLogger) (*Server, error) {
	s := &Server{
		Mux:          mux.NewRouter(),
		SessionStore: sessions.NewCookieStore([]byte(config.SessionKey)),
		Logger:       logger,
		Sanitizer:    bluemonday.UGCPolicy(),
		Config:       config,
	}
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}

func (s *Server) Configure() error {
	_, err := store.NewStore(s.Config.DatabaseURL)
	if err != nil {
		return errors.Wrap(err, "configuring server")
	}

	return nil
}

func Start() error {
	config := NewConfig()

	zapLogger, err := zap.NewProduction()
	if err != nil {
		return errors.Wrap(err, "starting server")
	}

	defer func() {
		if err := zapLogger.Sync(); err != nil {
			log.Println(errors.Wrap(err, "starting server"))
		}
	}()

	sugaredLogger := zapLogger.Sugar()

	server, err := NewServer(config, sugaredLogger)
	if err != nil {
		return errors.Wrap(err, "starting server")
	}

	if err := server.Configure(); err != nil {
		return errors.Wrap(err, "starting server")
	}

	log.Println("starting server at", config.BindAddr)
	return http.ListenAndServe(config.BindAddr, server)
}



