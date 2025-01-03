package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Client represents a single connection to the websocket. For every connection, a Client is
// created. The group broadcasts message to the connection using the send channel.
type Client struct {
	id    string
	conn  *websocket.Conn
	group *Group
	send  chan []byte
}

func NewClient(c *websocket.Conn, g *Group) *Client {
	return &Client{
		id:    generateRandomID(),
		conn:  c,
		group: g,
		send:  make(chan []byte, 256),
	}
}

// The Client's Read method reads messages from the websocket connection and writes them to
// c.group.broadcast channel. When the connection is closed, it unregisters the client from
// c.group and exits out of the goroutine
func (c *Client) Read() {
	defer func() {
		c.group.unregister <- c
		c.conn.Close()
	}()

	for {
		_, m, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Print(err)
			}
			break
		}

		// This isn't necessary. It's just used to differentiate the clients sending the message
		message := []byte(fmt.Sprintf("%s - %s", c.id, string(m)))
		c.group.broadcast <- message
	}
}

// The Client's Write method reads messages from c.send channel and writes them to the websocket
// connection. When the c.send channel is closed, it closes c.conn and exits out of the goroutine
func (c *Client) Write() {
	defer func() {
		c.conn.Close()
	}()

	for {
		message, ok := <-c.send
		if !ok {
			return
		}

		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}

func generateRandomID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)[:6]
}
