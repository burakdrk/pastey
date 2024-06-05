// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type ClipboardEntry struct {
	ID            int64     `json:"id"`
	EntryID       uuid.UUID `json:"entry_id"`
	UserID        int64     `json:"user_id"`
	FromDeviceID  int64     `json:"from_device_id"`
	ToDeviceID    int64     `json:"to_device_id"`
	EncryptedData string    `json:"encrypted_data"`
	CreatedAt     time.Time `json:"created_at"`
}

type Device struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	DeviceName string    `json:"device_name"`
	PublicKey  string    `json:"public_key"`
	CreatedAt  time.Time `json:"created_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	IpAddress    string    `json:"ip_address"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	DeviceID     int64     `json:"device_id"`
}

type User struct {
	ID              int64     `json:"id"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"password_hash"`
	Ispremium       bool      `json:"ispremium"`
	Isemailverified bool      `json:"isemailverified"`
	CreatedAt       time.Time `json:"created_at"`
}
