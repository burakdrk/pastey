package db

import (
	"context"
	"testing"

	"github.com/burakdrk/pastey/pastey-api/util"
	"github.com/stretchr/testify/require"
)

func TestSaveCopy(t *testing.T) {
	store := NewStore(testDB)

	user, err := store.GetUserById(context.Background(), 1)
	if err != nil {
		user = createRandomUser(t)
	}

	device1, err := store.GetDeviceById(context.Background(), 1)
	if err != nil {
		device1 = createRandomDevice(t, user.ID)
	}
	device2, err := store.GetDeviceById(context.Background(), 2)
	if err != nil {
		device2 = createRandomDevice(t, user.ID)
	}

	errs := make(chan error)
	results := make(chan []ClipboardEntry)

	arg := SaveCopyParams{
		UserID:       user.ID,
		FromDeviceID: device1.ID,
		Copies: []struct {
			ToDeviceID    int64  `json:"to_device_id"`
			EncryptedData string `json:"encrypted_data"`
		}{
			{
				ToDeviceID:    device2.ID,
				EncryptedData: util.RandomString(10),
			},
			{
				ToDeviceID:    device1.ID,
				EncryptedData: util.RandomString(10),
			},
		},
	}

	// run n concurrent transfer transaction
	n := 5
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.SaveCopy(context.Background(), arg)

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		entries := <-results
		require.Len(t, entries, len(arg.Copies))

		for i := range entries {
			require.NotEmpty(t, entries[i])
			require.Equal(t, arg.UserID, entries[i].UserID)
			require.Equal(t, arg.FromDeviceID, entries[i].FromDeviceID)
			require.Equal(t, arg.Copies[i].ToDeviceID, entries[i].ToDeviceID)
			require.Equal(t, arg.Copies[i].EncryptedData, entries[i].EncryptedData)
			require.NotZero(t, entries[i].CreatedAt)
			require.NotZero(t, entries[i].EntryID)
		}
	}
}
