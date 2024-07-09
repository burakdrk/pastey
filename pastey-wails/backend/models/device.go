package models

import "time"

type CreateDeviceResponse struct {
	Device                Device    `json:"device"`
	SessionID             string    `json:"session_id"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type Device struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	DeviceName string    `json:"device_name"`
	PublicKey  string    `json:"public_key"`
	CreatedAt  time.Time `json:"created_at"`
}
