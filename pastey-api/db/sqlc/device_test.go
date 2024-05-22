package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/burakdrk/pastey/pastey-api/util"
	"github.com/stretchr/testify/require"
)

func createRandomDevice(t *testing.T, userID int64) Device {
	if userID < 0 {
		user := createRandomUser(t)
		userID = user.ID
	}

	arg := CreateDeviceParams{
		UserID:     userID,
		DeviceName: util.RandomString(10),
		PublicKey:  util.RandomString(10),
	}

	device, err := testQueries.CreateDevice(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, device)

	require.Equal(t, arg.UserID, device.UserID)
	require.Equal(t, arg.DeviceName, device.DeviceName)
	require.Equal(t, arg.PublicKey, device.PublicKey)

	require.NotZero(t, device.CreatedAt)
	require.NotZero(t, device.ID)

	return device
}

func TestCreateDevice(t *testing.T) {
	createRandomDevice(t, -1)
}

func TestGetDeviceById(t *testing.T) {
	device1 := createRandomDevice(t, -1)

	device2, err := testQueries.GetDeviceById(context.Background(), device1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, device2)

	require.Equal(t, device1.UserID, device2.UserID)
	require.Equal(t, device1.DeviceName, device2.DeviceName)
	require.Equal(t, device1.PublicKey, device2.PublicKey)
	require.Equal(t, device1.CreatedAt, device2.CreatedAt)
	require.Equal(t, device1.ID, device2.ID)
}

func TestListUserDevices(t *testing.T) {
	user := createRandomUser(t)

	n := 4
	for i := 0; i < n; i++ {
		createRandomDevice(t, user.ID)
	}

	devices, err := testQueries.ListUserDevices(context.Background(), user.ID)
	require.NoError(t, err)
	require.Len(t, devices, n)

	for _, device := range devices {
		require.NotEmpty(t, device)
	}
}

func TestDeleteDevice(t *testing.T) {
	device1 := createRandomDevice(t, -1)

	err := testQueries.DeleteDevice(context.Background(), device1.ID)
	require.NoError(t, err)

	device2, err := testQueries.GetDeviceById(context.Background(), device1.ID)
	require.Error(t, err)
	require.Empty(t, device2)
}

func TestUpdateDevice(t *testing.T) {
	device1 := createRandomDevice(t, -1)

	arg := UpdateDeviceParams{
		ID:         device1.ID,
		DeviceName: sql.NullString{String: util.RandomString(10), Valid: true},
		PublicKey:  sql.NullString{String: util.RandomString(10), Valid: true},
	}

	device2, err := testQueries.UpdateDevice(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, device2)

	require.Equal(t, arg.ID, device2.ID)
	require.Equal(t, arg.DeviceName.String, device2.DeviceName)
	require.Equal(t, arg.PublicKey.String, device2.PublicKey)
	require.Equal(t, device1.UserID, device2.UserID)
	require.NotEqual(t, device1.DeviceName, device2.DeviceName)
	require.NotEqual(t, device1.PublicKey, device2.PublicKey)
}
