package models

import "time"

type Entry struct {
	ID             int64     `json:"id"`
	EntryID        string    `json:"entry_id"`
	UserID         int64     `json:"user_id"`
	FromDeviceID   int64     `json:"from_device_id"`
	ToDeviceID     int64     `json:"to_device_id"`
	EncryptedData  string    `json:"encrypted_data"`
	CreatedAt      time.Time `json:"created_at"`
	FromDeviceName string    `json:"from_device_name"`
}

type CopyRequest struct {
	FromDeviceID int64  `json:"from_device_id"`
	Copies       []Copy `json:"copies"`
}

type Copy struct {
	ToDeviceID    int64  `json:"to_device_id"`
	EncryptedData string `json:"encrypted_data"`
}
