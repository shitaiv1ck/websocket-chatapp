package main

import (
	"context"
	"os/signal"
	"syscall"

	core_logger "github.com/shitaiv1ck/realtime-chat/internal/core/logger"
	core_http_server "github.com/shitaiv1ck/realtime-chat/internal/core/server/http"
	core_ws_server "github.com/shitaiv1ck/realtime-chat/internal/core/server/ws"
	core_postgres "github.com/shitaiv1ck/realtime-chat/internal/core/store/postgres"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	log, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		panic(err)
	}

	log.Debug("init postgres store...")
	postgresStore := core_postgres.NewStore(log)
	if err := postgresStore.Open(ctx); err != nil {
		panic(err)
	}

	log.Debug("init ws server...")
	wsServer := core_ws_server.NewServer(log)
	go wsServer.Run(ctx)

	log.Debug("init http server...")
	httpServer := core_http_server.NewServer(nil, log)
	if err := httpServer.Run(ctx); err != nil {
		panic(err)
	}
}
