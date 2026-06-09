package core_ws_server

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
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
		c.closeSocket()
		c.server.leave <- c
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
			return
		}

	}
}

func (c *Client) Write() {
	defer c.closeSocket()

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
