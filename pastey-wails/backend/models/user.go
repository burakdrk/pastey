package models

import "time"

type LoginResponse struct {
	SessionID             string    `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  User      `json:"user"`
}

type RefreshTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

type User struct {
	ID              int64     `json:"id"`
	Email           string    `json:"email"`
	Ispremium       bool      `json:"ispremium"`
	Isemailverified bool      `json:"isemailverified"`
	CreatedAt       time.Time `json:"created_at"`
}
