// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateDevice(ctx context.Context, arg CreateDeviceParams) (Device, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (ClipboardEntry, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteDevice(ctx context.Context, id int64) error
	DeleteEntry(ctx context.Context, entryID uuid.UUID) error
	DeleteSession(ctx context.Context, id uuid.UUID) error
	DeleteSessionsByDevice(ctx context.Context, deviceID int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetDeviceById(ctx context.Context, id int64) (Device, error)
	GetEntriesForDevice(ctx context.Context, toDeviceID int64) ([]GetEntriesForDeviceRow, error)
	GetEntryByEntryId(ctx context.Context, entryID uuid.UUID) ([]ClipboardEntry, error)
	GetEntryByUser(ctx context.Context, userID int64) ([]ClipboardEntry, error)
	GetEntryByUserForUpdate(ctx context.Context, userID int64) ([]ClipboardEntry, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id int64) (User, error)
	ListUserDevices(ctx context.Context, userID int64) ([]Device, error)
	UpdateDevice(ctx context.Context, arg UpdateDeviceParams) (Device, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
