package db

import (
	"context"
	"testing"

	"github.com/burakdrk/pastey/pastey-api/util"
	"github.com/stretchr/testify/require"
)

func createRandomDevice(t *testing.T, userID int64) Device {
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
