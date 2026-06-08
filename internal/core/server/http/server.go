package core_http_server

import (
	"context"
	"errors"
	"net/http"

	core_logger "github.com/shitaiv1ck/realtime-chat/internal/core/logger"
	"go.uber.org/zap"
)

type Server struct {
	config Config
	router http.Handler
	log    *core_logger.Logger
}

func NewServer(router http.Handler, log *core_logger.Logger) *Server {
	return &Server{
		config: NewConfigMust(),
		router: router,
		log:    log,
	}
}

func (s *Server) Run(ctx context.Context) error {
	httpServer := &http.Server{
		Addr:    s.config.Addr,
		Handler: s.router,
	}

	errRun := make(chan error)
	go func() {
		defer close(errRun)

		s.log.Debug("start http server")
		err := httpServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			errRun <- err
		}
	}()

	select {
	case err := <-errRun:
		s.log.Error("http server's work", zap.Error(err))
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			httpServer.Close()
			s.log.Debug("http server is closed hard")
		}

		s.log.Debug("http server is closed correctly")
	}

	return nil
}
