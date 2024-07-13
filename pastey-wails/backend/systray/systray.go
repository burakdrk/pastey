package systray

import (
	"context"
	"runtime"

	"fyne.io/systray"
	"fyne.io/systray/example/icon"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Systray struct {
	ctx context.Context
}

func NewSystray() *Systray {
	return &Systray{}
}

func (s *Systray) Run(ctx context.Context) {
	s.ctx = ctx
	runtime.LockOSThread()

	systray.Run(s.onReady, s.onExit)
}

func (s *Systray) onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTooltip("Pastey")

	mOpen := systray.AddMenuItem("Open", "Open the app")
	mHide := systray.AddMenuItem("Hide", "Hide the app")
	mQuit := systray.AddMenuItem("Quit", "Quit the app")

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				s.onExit()
			case <-mHide.ClickedCh:
				wailsRuntime.Hide(s.ctx)
			case <-mOpen.ClickedCh:
				wailsRuntime.Show(s.ctx)
			}
		}
	}()
}

func (s *Systray) onExit() {
	wailsRuntime.Quit(s.ctx)
}
