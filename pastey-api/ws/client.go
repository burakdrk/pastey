package ws

import (
	"github.com/gorilla/websocket"
)

type Message struct {
	UserID        int64  `json:"user_id"`
	FromDeviceID  int64  `json:"from_device_id"`
	ToDeviceID    int64  `json:"to_device_id"`
	EncryptedData string `json:"encrypted_data"`
}

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	UserID   int64
	DeviceID int64
}

func NewClient(conn *websocket.Conn, userID, deviceID int64) *Client {
	return &Client{
		Conn:     conn,
		Message:  make(chan *Message, 256),
		UserID:   userID,
		DeviceID: deviceID,
	}
}

func (c *Client) Read(h *Hub) {
	defer func() {
		h.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Conn.WriteJSON(message)
		}
	}
}
