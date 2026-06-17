package core_ws_server

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	id        int
	socket    *websocket.Conn
	server    *Server
	receive   chan []byte
	closeOnce sync.Once
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func NewClient(id int, socket *websocket.Conn, server *Server) *Client {
	return &Client{
		id:        id,
		socket:    socket,
		server:    server,
		receive:   make(chan []byte, 256),
		closeOnce: sync.Once{},
	}
}

func (c *Client) closeSocket() {
	c.closeOnce.Do(func() {
		c.socket.Close()
	})
}

func (c *Client) Read() {
	defer func() {
		c.server.log.Debug("close read...")
		c.server.leave <- c
		c.closeSocket()
		c.server.log.Debug("read is closed")
	}()

	c.socket.SetReadLimit(maxMessageSize)
	c.socket.SetReadDeadline(time.Now().Add(pongWait))
	c.socket.SetPongHandler(func(string) error {
		c.socket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.socket.ReadMessage()
		if err != nil {
			c.server.log.Debug("pong", zap.Error(err))
			return
		}

	}
}

func (c *Client) Write() {
	defer func() {
		c.server.log.Debug("close write...")
		c.closeSocket()
		c.server.log.Debug("write is closed")
	}()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.receive:
			c.socket.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
				c.server.log.Debug("ping", zap.Error(err))
				return
			}
		case <-ticker.C:
			c.socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.socket.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
