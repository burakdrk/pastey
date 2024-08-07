package backend

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/burakdrk/pastey/pastey-wails/backend/clipboard"
	"github.com/burakdrk/pastey/pastey-wails/backend/crypto"
	"github.com/burakdrk/pastey/pastey-wails/backend/models"
	"github.com/burakdrk/pastey/pastey-wails/backend/storage"
	"github.com/burakdrk/pastey/pastey-wails/backend/systray"
	"github.com/burakdrk/pastey/pastey-wails/backend/utils"
	"github.com/go-resty/resty/v2"
)

// App struct
type App struct {
	ctx        context.Context
	api        *resty.Client
	Logger     *storage.Logger
	storage    storage.Storage
	clipboard  *clipboard.Clipboard
	ws         *clipboard.WSClient
	systray    *systray.Systray
	isLoggedIn bool
	deviceId   int64
}

// NewApp creates a new App application struct
func NewApp() *App {
	client := resty.New()
	client.SetBaseURL(utils.BASE_URL)
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
	ws := clipboard.NewWSClient()
	systray := systray.NewSystray()

	clipboard, err := clipboard.NewClipboard()
	if err != nil {
		logger.Log(err.Error())
		panic(err)
	}

	return &App{
		api:        client,
		Logger:     logger,
		storage:    storage,
		clipboard:  clipboard,
		systray:    systray,
		ws:         ws,
		isLoggedIn: false,
		deviceId:   0,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	a.api.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		if req.URL == "/token/refresh" {
			return nil
		}

		token, err := a.getAccessToken()
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	})

	a.clipboard.Listen(ctx, a.copyClipboard)
	go a.systray.Run(ctx)

	expiry, err := a.storage.Get("refresh_token_expiry")
	if err != nil {
		return
	}

	parsedExpiry, err := time.Parse(utils.DATE_FORMAT, expiry)
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
	a.ConnectToWS()
}

func (a *App) Shutdown(ctx context.Context) {
	a.storage.Close()
	a.ws.Close()
}

func (a *App) GetIsLoggedIn() bool {
	return a.isLoggedIn
}

func (a *App) GetDeviceId() int64 {
	return a.deviceId
}

func (a *App) GetConnectionStatus() bool {
	return a.ws.GetIsConnected()
}

type GetEntriesResponse struct {
	Entries []models.Entry `json:"entries"`
	Error   models.Error   `json:"error"`
}

func (a *App) GetEntries() GetEntriesResponse {
	var entries []models.Entry
	var errResp models.Error

	deviceId, err := a.storage.Get("device_id")
	if err != nil {
		a.Logger.Log(err.Error())
		return GetEntriesResponse{entries, models.GetDefaultError()}
	}

	res, err := a.api.R().
		SetError(&errResp).
		SetResult(&entries).
		Get(fmt.Sprintf("/devices/%s/entries", deviceId))
	if err != nil {
		a.Logger.Log(err.Error())
		return GetEntriesResponse{entries, models.GetDefaultError()}
	}

	if res.IsError() {
		return GetEntriesResponse{entries, errResp}
	}

	privateKey, err := a.storage.Get("private_key")
	if err != nil {
		a.Logger.Log(err.Error())
		return GetEntriesResponse{entries, models.GetDefaultError()}
	}

	for i := range entries {
		decrypted, err := crypto.DecryptData(entries[i].EncryptedData, privateKey)
		if err != nil {
			a.Logger.Log(err.Error())
			return GetEntriesResponse{entries, models.GetDefaultError()}
		}

		entries[i].EncryptedData = decrypted
	}

	return GetEntriesResponse{entries, models.Error{}}
}

func (a *App) GetDevices() ([]models.Device, models.Error) {
	var devices []models.Device
	var errResp models.Error

	res, err := a.api.R().
		SetError(&errResp).
		SetResult(&devices).
		Get("/devices")
	if err != nil {
		a.Logger.Log(err.Error())
		return devices, models.GetDefaultError()
	}

	if res.IsError() {
		return devices, errResp
	}

	return devices, models.Error{}
}

func (a *App) copyClipboard(data string) models.Error {
	var errResp models.Error

	devices, deviceErr := a.GetDevices()
	if deviceErr.Message != "" {
		return deviceErr
	}

	copies := []models.Copy{}
	for _, device := range devices {
		encrypted, err := crypto.EncryptData(data, device.PublicKey)
		if err != nil {
			a.Logger.Log(err.Error())
			return models.GetDefaultError()
		}

		copies = append(copies, models.Copy{
			ToDeviceID:    device.ID,
			EncryptedData: encrypted,
		})
	}

	req := models.CopyRequest{
		FromDeviceID: a.deviceId,
		Copies:       copies,
	}

	res, err := a.api.R().
		SetBody(req).
		SetError(&errResp).
		Post("/copy")
	if err != nil {
		a.Logger.Log(err.Error())
		return models.GetDefaultError()
	}

	if res.IsError() {
		return errResp
	}

	return models.Error{}
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
		a.deviceId = deviceResp.Device.ID
	}

	a.isLoggedIn = true
	return models.Error{}
}

func (a *App) DeleteEntry(entryId string) models.Error {
	var errResp models.Error

	res, err := a.api.R().
		SetError(&errResp).
		Delete(fmt.Sprintf("/entries/%s", entryId))
	if err != nil {
		a.Logger.Log(err.Error())
		return models.GetDefaultError()
	}

	if res.IsError() {
		return errResp
	}

	return models.Error{}
}

func (a *App) getAccessToken() (string, error) {
	token, err := a.storage.Get("access_token")
	if err != nil {
		return "", err
	}

	expiry, err := a.storage.Get("access_token_expiry")
	if err != nil {
		return "", err
	}

	parsedExpiry, err := time.Parse(utils.DATE_FORMAT, expiry)
	if err != nil {
		return "", err
	}

	if time.Now().UTC().After(parsedExpiry) {
		refreshToken, err := a.storage.Get("refresh_token")
		if err != nil {
			return "", err
		}

		var errResp models.Error
		var refreshResp models.RefreshTokenResponse

		res, err := a.api.R().
			SetBody(map[string]string{
				"refresh_token": refreshToken,
			}).
			SetError(&errResp).
			SetResult(&refreshResp).
			Post("/token/refresh")
		if err != nil {
			return "", err
		}

		if res.IsError() {
			return "", fmt.Errorf("error refreshing token: %s", errResp.Message)
		}

		a.storage.Save("access_token", refreshResp.AccessToken)
		a.storage.Save("access_token_expiry", refreshResp.AccessTokenExpiresAt.String())
		token = refreshResp.AccessToken
	}

	return token, nil
}

func (a *App) ConnectToWS() models.Error {
	token, err := a.getAccessToken()
	if err != nil {
		return models.Error{Message: err.Error()}
	}

	err = a.ws.Connect(a.ctx, fmt.Sprintf("%s?device_id=%d", utils.WEB_SOCKET_URL, a.deviceId), token)
	if err != nil {
		a.Logger.Log(err.Error())
		return models.Error{Message: err.Error()}
	}

	privateKey, err := a.storage.Get("private_key")
	if err != nil {
		a.Logger.Log(err.Error())
		return models.Error{Message: err.Error()}
	}

	go a.ws.Listen(a.clipboard, privateKey)
	return models.Error{}
}
