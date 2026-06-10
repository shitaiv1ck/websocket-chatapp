package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	core_logger "github.com/shitaiv1ck/realtime-chat/internal/core/logger"
	core_http_server "github.com/shitaiv1ck/realtime-chat/internal/core/server/http"
	core_ws_server "github.com/shitaiv1ck/realtime-chat/internal/core/server/ws"
	core_postgres "github.com/shitaiv1ck/realtime-chat/internal/core/store/postgres"
	core_middleware "github.com/shitaiv1ck/realtime-chat/internal/core/transport/middleware"
	sessions_repository "github.com/shitaiv1ck/realtime-chat/internal/features/sessions/repository"
	sessions_service "github.com/shitaiv1ck/realtime-chat/internal/features/sessions/service"
	sessions_http_transport "github.com/shitaiv1ck/realtime-chat/internal/features/sessions/transport/http"
	users_repository "github.com/shitaiv1ck/realtime-chat/internal/features/users/repository"
	users_service "github.com/shitaiv1ck/realtime-chat/internal/features/users/service"
	users_http_transport "github.com/shitaiv1ck/realtime-chat/internal/features/users/transport/http"
	users_ws_transport "github.com/shitaiv1ck/realtime-chat/internal/features/users/transport/ws"
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

	log.Debug("init feature: users...")
	usersRep := users_repository.NewRepository(postgresStore)
	usersWS := users_ws_transport.NewWSTransport(wsServer)
	usersService := users_service.NewService(usersRep, usersWS)
	usersHTTP := users_http_transport.NewHTTPTransport(usersService)

	log.Debug("init feature: sessions...")
	sessionsRep := sessions_repository.NewRepository(postgresStore)
	sessionsService := sessions_service.NewService(sessionsRep, usersRep)
	sessionsHTTP := sessions_http_transport.NewHTTPTransport(sessionsService)

	protected := http.NewServeMux()
	protected.Handle("GET /users/me", usersHTTP.GetMeHandler())
	protected.Handle("PATCH /users", usersHTTP.PatchUserHandler())
	protected.Handle("DELETE /sessions", sessionsHTTP.DeleteSessionHandler())

	protectedHandler := core_middleware.ProtectedMiddleware(protected, sessionsService)

	common := http.NewServeMux()
	common.Handle("/ws", wsServer)
	common.Handle("POST /api/users", usersHTTP.CreateUserHandler())
	common.Handle("GET /api/users", usersHTTP.GetUsersHandler())
	common.Handle("POST /api/sessions", sessionsHTTP.CreateSessionHandler())
	common.Handle("/api/protected/", http.StripPrefix("/api/protected", protectedHandler))

	commonHandler := core_middleware.CommonMiddleware(common, log)

	log.Debug("init http server...")
	httpServer := core_http_server.NewServer(commonHandler, log)
	if err := httpServer.Run(ctx); err != nil {
		panic(err)
	}
}
