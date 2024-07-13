package clipboard

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/burakdrk/pastey/pastey-wails/backend/crypto"
	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// WSClient struct
type WSClient struct {
	conn        *websocket.Conn
	ctx         context.Context
	ticker      *time.Ticker
	isConnected bool
	mu          sync.Mutex
}

type WSEntry struct {
	UserID        int64  `json:"user_id"`
	FromDeviceID  int64  `json:"from_device_id"`
	ToDeviceID    int64  `json:"to_device_id"`
	EncryptedData string `json:"encrypted_data"`
}

// NewWSClient creates a new WSClient struct
func NewWSClient() *WSClient {
	return &WSClient{
		isConnected: false,
	}
}

func (c *WSClient) Connect(ctx context.Context, url string, token string) error {
	conn, _, err := websocket.DefaultDialer.Dial(url, http.Header{
		"Authorization": []string{"Bearer " + token},
	})
	if err != nil {
		return err
	}

	c.setIsConnected(true)

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	conn.SetPongHandler(func(_ string) error {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		return nil
	})

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					runtime.EventsEmit(ctx, "ws:disconnected", nil)
					c.setIsConnected(false)
					return
				}
			}
		}
	}()

	c.conn = conn
	c.ctx = ctx
	c.ticker = ticker
	return nil
}

// Run this before closing the app
func (c *WSClient) Close() error {
	c.ticker.Stop()
	return c.conn.Close()
}

// Blocking, call this in a goroutine
func (c *WSClient) Listen(clipboard *Clipboard, privateKey string) {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			runtime.EventsEmit(c.ctx, "ws:disconnected", nil)
			c.setIsConnected(false)
			break
		}

		var entry WSEntry

		err = json.Unmarshal(message, &entry)
		if (err != nil) || (entry.ToDeviceID == entry.FromDeviceID) {
			continue
		}

		decrypted, err := crypto.DecryptData(entry.EncryptedData, privateKey)
		if err != nil {
			fmt.Println(err)
			continue
		}

		clipboard.Set(decrypted)
		entry.EncryptedData = decrypted
		runtime.EventsEmit(c.ctx, "ws:entry", entry)
	}
}

func (c *WSClient) setIsConnected(isConnected bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.isConnected = isConnected
}

func (c *WSClient) GetIsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.isConnected
}
