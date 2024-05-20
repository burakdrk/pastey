// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	AcquireAdvisoryLock(ctx context.Context, pgAdvisoryLock int64) error
	CreateDevice(ctx context.Context, arg CreateDeviceParams) (Device, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (ClipboardEntry, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteDevice(ctx context.Context, id int64) error
	DeleteEntry(ctx context.Context, entryID uuid.UUID) error
	DeleteUser(ctx context.Context, id int64) error
	GetDeviceById(ctx context.Context, id int64) (Device, error)
	GetEntriesForDevice(ctx context.Context, toDeviceID int64) ([]ClipboardEntry, error)
	GetEntryByEntryId(ctx context.Context, entryID uuid.UUID) ([]ClipboardEntry, error)
	GetEntryByUserForUpdate(ctx context.Context, userID int64) ([]ClipboardEntry, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id int64) (User, error)
	ListUserDevices(ctx context.Context, userID int64) ([]Device, error)
	ReleaseAdvisoryLock(ctx context.Context, pgAdvisoryUnlock int64) error
	UpdateDevice(ctx context.Context, arg UpdateDeviceParams) (Device, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
