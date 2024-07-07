package backend

import (
	"context"
	"fmt"
	"net/http"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Shutdown(ctx context.Context) {
	a.ctx.Done()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Login(email string, password string) {\
	reqBody := map[string]string{
		"email":    email,
		"password": password,
	}

	http.Post("https://api.burakduruk.com/v1/users/login", "application/json", reqBody)
}
