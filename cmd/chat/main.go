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
	chats_repository "github.com/shitaiv1ck/realtime-chat/internal/features/chats/repository"
	chats_service "github.com/shitaiv1ck/realtime-chat/internal/features/chats/service"
	chats_http_transport "github.com/shitaiv1ck/realtime-chat/internal/features/chats/transport/http"
	chats_ws_transport "github.com/shitaiv1ck/realtime-chat/internal/features/chats/transport/ws"
	friendrequests_respository "github.com/shitaiv1ck/realtime-chat/internal/features/friendrequests/respository"
	friendrequests_service "github.com/shitaiv1ck/realtime-chat/internal/features/friendrequests/service"
	friendrequests_http_transport "github.com/shitaiv1ck/realtime-chat/internal/features/friendrequests/transport/http"
	friendrequests_ws_transport "github.com/shitaiv1ck/realtime-chat/internal/features/friendrequests/transport/ws"
	friendships_repository "github.com/shitaiv1ck/realtime-chat/internal/features/friendships/repository"
	friendships_service "github.com/shitaiv1ck/realtime-chat/internal/features/friendships/service"
	friendships_http_transport "github.com/shitaiv1ck/realtime-chat/internal/features/friendships/transport/http"
	friendships_ws_transport "github.com/shitaiv1ck/realtime-chat/internal/features/friendships/transport/ws"
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
	postgresStore, err := core_postgres.NewConnPool(ctx, core_postgres.NewConfigMust())
	if err != nil {
		panic(err)
	}
	defer postgresStore.Close()

	log.Debug("init ws server...")
	wsServer := core_ws_server.NewServer(log)
	go wsServer.Run(ctx)

	log.Debug("init repositories...")
	usersRep := users_repository.NewRepository(postgresStore)
	sessionsRep := sessions_repository.NewRepository(postgresStore)
	friendRequestsRep := friendrequests_respository.NewRepository(postgresStore)
	friendshipsRep := friendships_repository.NewRepository(postgresStore)
	chatsRepository := chats_repository.NewRepository(postgresStore)

	log.Debug("init ws transports...")
	usersWS := users_ws_transport.NewWSTransport(wsServer)
	friendRequestsWS := friendrequests_ws_transport.NewWSTransport(wsServer)
	friendshipsWS := friendships_ws_transport.NewWSTransport(wsServer)
	chatsWS := chats_ws_transport.NewWSTransport(wsServer)

	log.Debug("init services...")
	usersService := users_service.NewService(usersRep, usersWS)
	sessionsService := sessions_service.NewService(sessionsRep, usersRep)
	friendRequestsService := friendrequests_service.NewService(friendRequestsRep, friendshipsRep, friendRequestsWS)
	friendshipsService := friendships_service.NewService(friendshipsRep, friendRequestsRep, friendshipsWS)
	chatsService := chats_service.NewService(chatsRepository, friendshipsRep, chatsWS)

	log.Debug("init http transports...")
	usersHTTP := users_http_transport.NewHTTPTransport(usersService)
	sessionsHTTP := sessions_http_transport.NewHTTPTransport(sessionsService)
	friendRequestsHTTP := friendrequests_http_transport.NewTransport(friendRequestsService)
	friendshipsHTTP := friendships_http_transport.NewHTTPTransport(friendshipsService)
	chatsHTTP := chats_http_transport.NewHTTPTransport(chatsService)

	protected := http.NewServeMux()
	protected.Handle("GET /users/me", usersHTTP.GetMeHandler())
	protected.Handle("PATCH /users", usersHTTP.PatchUserHandler())
	protected.Handle("DELETE /sessions", sessionsHTTP.DeleteSessionHandler())
	protected.Handle("POST /friend-requests", friendRequestsHTTP.CreateFriendRequestHandler())
	protected.Handle("GET /friend-requests", friendRequestsHTTP.GetFriendRequestsHandler())
	protected.Handle("DELETE /friend-requests/{friend_request_id}", friendRequestsHTTP.DeleteFriendRequestHandler())
	protected.Handle("POST /friendships", friendshipsHTTP.CreateFriendshipHandler())
	protected.Handle("GET /friendships", friendshipsHTTP.GetFriendshipsHandler())
	protected.Handle("DELETE /friendships/{friendship_id}", friendshipsHTTP.DeleteFriendshipHandler())
	protected.Handle("POST /chats", chatsHTTP.CreateOrGetChatHandler())
	protected.Handle("GET /chats", chatsHTTP.GetChatsHandler())
	protected.Handle("DELETE /chats/{chat_id}", chatsHTTP.DeleteChatHandler())

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
