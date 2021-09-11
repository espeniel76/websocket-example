package pkgs

import (
	"log"

	"github.com/gorilla/websocket"
)

const (
	maxMessageSize = 1024
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	World *World
	Conn  *websocket.Conn
	Send  chan []byte
}

func NewClient(w *World, c *websocket.Conn) (client *Client) {
	client = &Client{
		World: w,
		Conn:  c,
		Send:  make(chan []byte, 256),
	}
	go client.readPump()
	go client.writePump()
	return client
}

func (c *Client) readPump() {
	defer func() {
		c.World.ChanLeave <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.World.broadcast <- message
	}
}

func (c *Client) writePump() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		message, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		c.Conn.WriteMessage(websocket.TextMessage, message)
	}
	// for {
	// 	message, ok := <-c.Send:
	// 		c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	// 		if !ok {
	// 			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
	// 			return
	// 		}
	// 		c.Conn.WriteMessage(websocket.TextMessage, message)
	// select {
	// case message, ok := <-c.Send:
	// 	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	// 	if !ok {
	// 		c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
	// 		return
	// 	}
	// 	c.Conn.WriteMessage(websocket.TextMessage, message)
	// case <-ticker.C:
	// 	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	// 	if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
	// 		return
	// 	}
	// }
	// }
}
