package clipboard

import (
	"context"

	"github.com/burakdrk/pastey/pastey-wails/backend/models"
	"golang.design/x/clipboard"
)

type Clipboard struct {
}

func NewClipboard() (*Clipboard, error) {
	err := clipboard.Init()
	if err != nil {
		return nil, err
	}

	return &Clipboard{}, nil
}

func (c *Clipboard) Start(ctx context.Context, callback func(string) models.Error) {
	ch := clipboard.Watch(ctx, clipboard.FmtText)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-ch:
				go callback(string(data))
			}
		}
	}()
}

func (c *Clipboard) Set(data string) {
	clipboard.Write(clipboard.FmtText, []byte(data))
}

func (c *Clipboard) Get() string {
	return string(clipboard.Read(clipboard.FmtText))
}
