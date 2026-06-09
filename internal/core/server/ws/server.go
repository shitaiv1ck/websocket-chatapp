package core_ws_server

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
	core_logger "github.com/shitaiv1ck/realtime-chat/internal/core/logger"
	core_repsponse "github.com/shitaiv1ck/realtime-chat/internal/core/transport/repsponse"
	core_utils "github.com/shitaiv1ck/realtime-chat/internal/core/utils"
	"go.uber.org/zap"
)

type Server struct {
	clients  map[int]*Client
	join     chan *Client
	leave    chan *Client
	upgrader *websocket.Upgrader
	mtx      sync.RWMutex
	log      *core_logger.Logger
}

const (
	readBufferSize  = 1024
	writeBufferSize = 1024
)

func NewServer(logger *core_logger.Logger) *Server {
	return &Server{
		clients: map[int]*Client{},
		join:    make(chan *Client),
		leave:   make(chan *Client),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  readBufferSize,
			WriteBufferSize: writeBufferSize,
		},
		mtx: sync.RWMutex{},
		log: logger,
	}
}

func (s *Server) Run(ctx context.Context) {
	s.log.Debug("start ws server")
	for {
		select {
		case <-ctx.Done():
			s.log.Debug("ws server is closed")
			return
		case client := <-s.join:
			s.mtx.Lock()
			if _, ok := s.clients[client.id]; !ok {
				s.clients[client.id] = client
				s.log.Debug("client join", zap.Int("client-id", client.id))
			}
			s.mtx.Unlock()
		case client := <-s.leave:
			s.mtx.Lock()
			if _, ok := s.clients[client.id]; ok {
				delete(s.clients, client.id)
				close(client.receive)
			}
			s.mtx.Unlock()
			s.log.Debug("client leave", zap.Int("client-id", client.id))
		}
	}
}

func (s *Server) Broadcast(msg []byte) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	for _, client := range s.clients {
		select {
		case client.receive <- msg:
		default:
		}
	}
}

func (s *Server) NotifyClient(id int, msg []byte) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	client, ok := s.clients[id]
	if !ok {
		return
	}

	select {
	case client.receive <- msg:
	default:
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	responseHandler := core_repsponse.NewResponseWriter(w)

	userID, err := core_utils.GetIntFromContext(req.Context(), "user_id")
	if err != nil {
		responseHandler.ErrorResponse(core_errors.ErrCoockie, "failed to authentication")

		return
	}

	socket, err := s.upgrader.Upgrade(w, req, nil)
	if err != nil {
		s.log.Error("failed to upgrade http request to websocket", zap.Error(err))
		return
	}

	client := NewClient(*userID, socket, s)

	s.join <- client

	go client.Read()
	client.Write()
}
