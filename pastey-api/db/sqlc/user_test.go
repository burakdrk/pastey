package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/burakdrk/pastey/pastey-api/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(10))
	require.NoError(t, err)

	arg := CreateUserParams{
		Email:        util.RandomString(7) + "@gmail.com",
		PasswordHash: hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)

	require.False(t, user.Ispremium)

	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.ID)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByEmail(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUserByEmail(context.Background(), user1.Email)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
	require.Equal(t, user1.ID, user2.ID)
}

func TestGetUserById(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUserById(context.Background(), user1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
	require.Equal(t, user1.ID, user2.ID)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)

	hashedPassword, err := util.HashPassword(util.RandomString(10))
	require.NoError(t, err)

	arg := UpdateUserParams{
		Email: sql.NullString{
			String: util.RandomString(7) + "@gmail.com",
			Valid:  true,
		},
		PasswordHash: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
		Ispremium: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		ID: user1.ID,
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.NotEqual(t, user1.Email, user2.Email)
	require.Equal(t, arg.Email.String, user2.Email)

	require.NotEqual(t, user1.PasswordHash, user2.PasswordHash)
	require.Equal(t, arg.PasswordHash.String, user2.PasswordHash)

	require.NotEqual(t, user1.Ispremium, user2.Ispremium)
	require.Equal(t, arg.Ispremium.Bool, user2.Ispremium)

	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
	require.Equal(t, user1.ID, user2.ID)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user1.ID)

	require.NoError(t, err)

	user2, err := testQueries.GetUserById(context.Background(), user1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}
