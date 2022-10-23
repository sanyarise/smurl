package server

import (
	"context"
	"net/http"
	"time"

	"github.com/sanyarise/smurl/internal/usecases/repos/smurlrepo"

	"go.uber.org/zap"
)

type Server struct {
	srv    http.Server
	logger *zap.Logger
}

func NewServer(addr string, h http.Handler, l *zap.Logger, rto int, wto int, rhto int) *Server {
	l.Debug("Enter in server func NewServer()")
	s := &Server{}

	s.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       time.Duration(rto) * time.Second,
		WriteTimeout:      time.Duration(wto) * time.Second,
		ReadHeaderTimeout: time.Duration(rhto) * time.Second,
	}
	return s
}

func (s *Server) Stop() {
	s.logger.Debug("Enter in server func Stop()")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	s.srv.Shutdown(ctx)
	cancel()
}

func (s *Server) Start(smr *smurlrepo.SmurlStorage) {
	s.logger.Debug("Enter in server func Start()")
	go s.srv.ListenAndServe()
}
