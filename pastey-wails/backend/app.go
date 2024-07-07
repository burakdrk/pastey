package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/burakdrk/pastey/pastey-wails/backend/models"
	"github.com/go-resty/resty/v2"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	client := resty.New()
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

func (a *App) Login(email string, password string) models.Error {
	reqBody, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	if err != nil {
		fmt.Println(err)
		return models.Error{Message: "Internal server error"}
	}

	res, err := http.Post("https://api.burakduruk.com/v1/users/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err)
		return models.Error{Message: "Internal server error"}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return models.Error{Message: "Internal server error"}
	}

	var response models.Error
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshaling response body:", err)
		return models.Error{Message: "Internal server error"}
	}

	fmt.Println("Response:", response)
	return response
}
