package utils

import (
	"fmt"
	"time"

	"github.com/burakdrk/pastey/pastey-wails/backend/models"
	"github.com/burakdrk/pastey/pastey-wails/backend/storage"
	"github.com/go-resty/resty/v2"
)

func AccessTokenMiddleware(c *resty.Client, req *resty.Request, storage storage.Storage) error {
	token, err := storage.Get("access_token")
	if err != nil {
		return nil
	}

	expiry, err := storage.Get("access_token_expiry")
	if err != nil {
		return err
	}

	parsedExpiry, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", expiry)
	if err != nil {
		return err
	}

	if time.Now().UTC().After(parsedExpiry) {
		refreshToken, err := storage.Get("refresh_token")
		if err != nil {
			return err
		}

		var errResp models.Error
		var refreshResp models.RefreshTokenResponse

		res, err := c.R().
			SetBody(map[string]string{
				"refresh_token": refreshToken,
			}).
			SetError(&errResp).
			SetResult(&refreshResp).
			Post("/token/refresh")
		if err != nil {
			return err
		}

		if res.IsError() {
			return fmt.Errorf("error refreshing token: %s", errResp.Message)
		}

		storage.Save("access_token", refreshResp.AccessToken)
		storage.Save("access_token_expiry", refreshResp.AccessTokenExpiresAt.String())
		token = refreshResp.AccessToken
	}

	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}
