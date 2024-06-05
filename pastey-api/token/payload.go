package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrExpiredToken = errors.New("token is expired")
var ErrInvalidToken = errors.New("token is invalid")

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewPayload(userID int64, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:        tokenId,
		UserID:    userID,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(duration),
	}, nil
}

func (p *Payload) Valid() error {
	if time.Now().UTC().After(p.ExpiresAt) {
		return ErrExpiredToken
	}

	return nil
}
