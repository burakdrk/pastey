package backend

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/burakdrk/pastey/pastey-wails/backend/crypto"
	"github.com/burakdrk/pastey/pastey-wails/backend/models"
	"github.com/burakdrk/pastey/pastey-wails/backend/storage"
	"github.com/burakdrk/pastey/pastey-wails/backend/utils"
	"github.com/go-resty/resty/v2"
)

// App struct
type App struct {
	ctx        context.Context
	api        *resty.Client
	Logger     *storage.Logger
	storage    storage.Storage
	isLoggedIn bool
	deviceId   int64
}

// NewApp creates a new App application struct
func NewApp() *App {
	client := resty.New()
	client.SetBaseURL("https://api.burakduruk.com/v1")
	client.SetHeader("Accept", "application/json")
	client.SetHeader("Content-Type", "application/json")

	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	appConfigPath := filepath.Join(configDir, "pastey")
	err = os.MkdirAll(appConfigPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	logger := storage.NewLogger(appConfigPath)
	storage := storage.NewSQLiteStorage(appConfigPath)

	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		return utils.AccessTokenMiddleware(c, req, storage)
	})

	return &App{
		api:        client,
		Logger:     logger,
		storage:    storage,
		isLoggedIn: false,
		deviceId:   0,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	expiry, err := a.storage.Get("refresh_token_expiry")
	if err != nil {
		return
	}

	parsedExpiry, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", expiry)
	if err != nil {
		a.Logger.Log(err.Error())
		return
	}

	if time.Now().UTC().Before(parsedExpiry) {
		a.isLoggedIn = true
	}

	did, err := a.storage.Get("device_id")
	if err != nil {
		return
	}

	deviceId, err := strconv.ParseInt(did, 10, 64)
	if err != nil {
		a.Logger.Log(err.Error())
		return
	}

	a.deviceId = deviceId
}

func (a *App) Shutdown(ctx context.Context) {
	a.ctx.Done()
	a.storage.Close()
}

func (a *App) GetIsLoggedIn() bool {
	return a.isLoggedIn
}

func (a *App) Login(email string, password string) models.Error {
	var errResp models.Error
	var loginResp models.LoginResponse

	var req map[string]interface{}

	if a.deviceId == 0 {
		req = map[string]interface{}{
			"email":    email,
			"password": password,
		}
	} else {
		req = map[string]interface{}{
			"email":     email,
			"password":  password,
			"device_id": a.deviceId,
		}
	}

	res, err := a.api.R().
		SetBody(req).
		SetError(&errResp).
		SetResult(&loginResp).
		Post("/users/login")
	if err != nil {
		a.Logger.Log(err.Error())
		return models.GetDefaultError()
	}

	if res.IsError() {
		return errResp
	}

	a.storage.Save("access_token", loginResp.AccessToken)
	a.storage.Save("access_token_expiry", loginResp.AccessTokenExpiresAt.String())
	if loginResp.RefreshToken != "" && !loginResp.RefreshTokenExpiresAt.IsZero() {
		a.storage.Save("refresh_token", loginResp.RefreshToken)
		a.storage.Save("refresh_token_expiry", loginResp.RefreshTokenExpiresAt.String())
	}

	if a.deviceId == 0 {
		var deviceResp models.CreateDeviceResponse

		priv, pub, keyErr := crypto.GenerateKeyPair(2048)
		if keyErr != nil {
			a.Logger.Log(keyErr.Error())
			return models.GetDefaultError()
		}

		a.storage.Save("private_key", priv)
		a.storage.Save("public_key", pub)

		hostname, err := os.Hostname()
		if err != nil {
			a.Logger.Log(err.Error())
			return models.GetDefaultError()
		}

		res, err := a.api.R().
			SetBody(map[string]string{
				"device_name": hostname,
				"public_key":  pub,
			}).
			SetError(&errResp).
			SetResult(&deviceResp).
			Post("/devices")
		if err != nil {
			a.Logger.Log(err.Error())
			return models.GetDefaultError()
		}

		if res.IsError() {
			return errResp
		}

		a.storage.Save("device_id", fmt.Sprint(deviceResp.Device.ID))
		a.storage.Save("refresh_token", deviceResp.RefreshToken)
		a.storage.Save("refresh_token_expiry", deviceResp.RefreshTokenExpiresAt.String())
	}

	return models.Error{}
}
