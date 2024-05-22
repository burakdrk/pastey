package db

import (
	"context"
	"testing"
	"time"

	"github.com/burakdrk/pastey/pastey-api/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) ClipboardEntry {
	user := createRandomUser(t)

	device1 := createRandomDevice(t, user.ID)
	device2 := createRandomDevice(t, user.ID)

	arg := CreateEntryParams{
		EntryID:       uuid.New(),
		UserID:        user.ID,
		CreatedAt:     time.Now().UTC(),
		FromDeviceID:  device1.ID,
		ToDeviceID:    device2.ID,
		EncryptedData: util.RandomString(10),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.UserID, entry.UserID)
	require.Equal(t, arg.EntryID, entry.EntryID)
	require.Equal(t, arg.FromDeviceID, entry.FromDeviceID)
	require.Equal(t, arg.ToDeviceID, entry.ToDeviceID)
	require.Equal(t, arg.EncryptedData, entry.EncryptedData)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntryByEntryId(t *testing.T) {
	entry1 := createRandomEntry(t)

	entries, err := testQueries.GetEntryByEntryId(context.Background(), entry1.EntryID)
	entry2 := entries[0]
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.UserID, entry2.UserID)
	require.Equal(t, entry1.EntryID, entry2.EntryID)
	require.Equal(t, entry1.FromDeviceID, entry2.FromDeviceID)
	require.Equal(t, entry1.ToDeviceID, entry2.ToDeviceID)
	require.Equal(t, entry1.EncryptedData, entry2.EncryptedData)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
	require.Equal(t, entry1.ID, entry2.ID)
}

func TestGetEntriesForDevice(t *testing.T) {
	entry1 := createRandomEntry(t)

	entries, err := testQueries.GetEntriesForDevice(context.Background(), entry1.ToDeviceID)
	entry2 := entries[0]
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.UserID, entry2.UserID)
	require.Equal(t, entry1.EntryID, entry2.EntryID)
	require.Equal(t, entry1.FromDeviceID, entry2.FromDeviceID)
	require.Equal(t, entry1.ToDeviceID, entry2.ToDeviceID)
	require.Equal(t, entry1.EncryptedData, entry2.EncryptedData)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
	require.Equal(t, entry1.ID, entry2.ID)
}

func TestGetEntryByUser(t *testing.T) {
	entry1 := createRandomEntry(t)

	entries, err := testQueries.GetEntryByUser(context.Background(), entry1.UserID)
	entry2 := entries[0]
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.UserID, entry2.UserID)
	require.Equal(t, entry1.EntryID, entry2.EntryID)
	require.Equal(t, entry1.FromDeviceID, entry2.FromDeviceID)
	require.Equal(t, entry1.ToDeviceID, entry2.ToDeviceID)
	require.Equal(t, entry1.EncryptedData, entry2.EncryptedData)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
	require.Equal(t, entry1.ID, entry2.ID)
}

func TestDeleteEntry(t *testing.T) {
	entry1 := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry1.EntryID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntryByEntryId(context.Background(), entry1.EntryID)
	require.NoError(t, err)
	require.Empty(t, entry2)
}
