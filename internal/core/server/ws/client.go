package core_ws_server

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	id        int
	socket    *websocket.Conn
	server    *Server
	receive   chan []byte
	closeOnce sync.Once
}

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

	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}

		c.server.MessageHandler(c, msg)
	}
}

func (c *Client) Write() {
	defer c.closeSocket()

	for msg := range c.receive {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}
