package core_postgres

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"

	_ "github.com/lib/pq"
	core_logger "github.com/shitaiv1ck/realtime-chat/internal/core/logger"
)

type Store struct {
	*sql.DB

	config Config
	log    *core_logger.Logger
}

func NewStore(logger *core_logger.Logger) *Store {
	return &Store{
		config: NewConfigMust(),
		log:    logger,
	}
}

func (s *Store) Open(ctx context.Context) error {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		s.config.User,
		s.config.Password,
		s.config.Host,
		s.config.Port,
		s.config.DB,
	)

	s.log.Debug("open conn to postgres store")
	db, err := sql.Open("postgres", url)
	if err != nil {
		s.log.Error("failed to open conn to postgres store", zap.Error(err))
		return err
	}

	if err := db.Ping(); err != nil {
		s.log.Error("failed to ping conn to postgres store", zap.Error(err))
		return err
	}

	s.DB = db

	go func() {
		<-ctx.Done()
		s.Close()
		s.log.Debug("conn to postgres store is closed")
	}()

	return nil
}
