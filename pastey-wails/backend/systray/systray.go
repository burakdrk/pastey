package systray

import (
	"fyne.io/systray"
	"fyne.io/systray/example/icon"
)

type Systray struct {
	quit chan struct{}
}

func NewSystray() *Systray {
	return &Systray{
		quit: make(chan struct{}),
	}
}

func (s *Systray) Run() {
	systray.Run(s.onReady, s.onExit)
}

func (s *Systray) onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTooltip("Pastey")

	mQuit := systray.AddMenuItem("Quit", "Quit Pastey")

	go func() {
		<-mQuit.ClickedCh
		s.quit <- struct{}{}
	}()
}

func (s *Systray) onExit() {
	close(s.quit)
	systray.Quit()
}
