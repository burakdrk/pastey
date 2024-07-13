package clipboard

import (
	"context"
	"sync"

	"github.com/burakdrk/pastey/pastey-wails/backend/models"
	"golang.design/x/clipboard"
)

type Clipboard struct {
	internalChange bool
	mu             sync.Mutex
}

func NewClipboard() (*Clipboard, error) {
	err := clipboard.Init()
	if err != nil {
		return nil, err
	}

	return &Clipboard{}, nil
}

func (c *Clipboard) Listen(ctx context.Context, callback func(string) models.Error) {
	ch := clipboard.Watch(ctx, clipboard.FmtText)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-ch:
				c.mu.Lock()
				if c.internalChange {
					c.internalChange = false
					c.mu.Unlock()
					continue
				}
				c.mu.Unlock()

				go callback(string(data))
			}
		}
	}()
}

func (c *Clipboard) Set(data string) {
	c.mu.Lock()
	c.internalChange = true
	c.mu.Unlock()

	clipboard.Write(clipboard.FmtText, []byte(data))
}

func (c *Clipboard) Get() string {
	return string(clipboard.Read(clipboard.FmtText))
}
