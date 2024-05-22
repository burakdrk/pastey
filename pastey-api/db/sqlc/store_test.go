package db

import (
	"context"
	"testing"

	"github.com/burakdrk/pastey/pastey-api/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSaveCopy1(t *testing.T) {
	store := NewStore(testDB)

	user, err := store.GetUserById(context.Background(), 1)
	if err != nil {
		user = createRandomUser(t)
	}

	user2, err := store.GetUserById(context.Background(), 2)
	if err != nil {
		user2 = createRandomUser(t)
	}

	device1, err := store.GetDeviceById(context.Background(), 1)
	if err != nil {
		device1 = createRandomDevice(t, user.ID)
	}
	device2, err := store.GetDeviceById(context.Background(), 2)
	if err != nil {
		device2 = createRandomDevice(t, user.ID)
	}
	device3, err := store.GetDeviceById(context.Background(), 3)
	if err != nil {
		device3 = createRandomDevice(t, user2.ID)
	}
	device4, err := store.GetDeviceById(context.Background(), 4)
	if err != nil {
		device4 = createRandomDevice(t, user2.ID)
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

	arg2 := SaveCopyParams{
		UserID:       user2.ID,
		FromDeviceID: device3.ID,
		Copies: []struct {
			ToDeviceID    int64  `json:"to_device_id"`
			EncryptedData string `json:"encrypted_data"`
		}{
			{
				ToDeviceID:    device4.ID,
				EncryptedData: util.RandomString(10),
			},
			{
				ToDeviceID:    device3.ID,
				EncryptedData: util.RandomString(10),
			},
		},
	}

	// run n concurrent transfer transaction
	n := 2
	for i := 0; i < n; i++ {
		go func(i int) {
			var argg SaveCopyParams
			if i == 0 {
				argg = arg
			} else {
				argg = arg2
			}
			result, err := store.SaveCopy(context.Background(), argg)
			t.Logf("Transaction %d: Error %v\n", i, err)
			for _, entry := range result {
				t.Logf("\tTransaction %d: %v\n", i, entry)
			}
			errs <- err
			results <- result
		}(i)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		entries := <-results
		// require.Len(t, entries, len(arg.Copies))

		for i := range entries {
			require.NotEmpty(t, entries[i])
			// require.Equal(t, arg.UserID, entries[i].UserID)
			// require.Equal(t, arg.FromDeviceID, entries[i].FromDeviceID)
			// require.Equal(t, arg.Copies[i].ToDeviceID, entries[i].ToDeviceID)
			// require.Equal(t, arg.Copies[i].EncryptedData, entries[i].EncryptedData)
			// require.NotZero(t, entries[i].CreatedAt)
			// require.NotZero(t, entries[i].EntryID)
		}
	}
}

func TestSaveCopy2(t *testing.T) {
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
		go func(i int) {
			result, err := store.SaveCopy(context.Background(), arg)
			t.Logf("Transaction %d: Error %v\n", i, err)
			for _, entry := range result {
				t.Logf("\tTransaction %d: %v\n", i, entry)
			}
			errs <- err
			results <- result
		}(i)
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

	newEntires, err := testQueries.GetEntryByUser(context.Background(), user.ID)
	require.NoError(t, err)
	count := 0

	seen := make(map[uuid.UUID]bool)
	for _, entry := range newEntires {
		if _, ok := seen[entry.EntryID]; ok {
			continue
		}
		seen[entry.EntryID] = true
		count++
	}

	require.Equal(t, 2, count)
}
