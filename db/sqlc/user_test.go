package db

import (
	"context"
	"testing"
	"time"

	"github.com/ahmadfarhanstwn/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username: util.RandomOwner(),
		Email: util.RandomEmail(),
		HashPassword: hashedPassword,
		FullName: util.RandomOwner(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	//make sure no error
	require.NoError(t, err)
	//make sure the function doesn't return empty value
	require.NotEmpty(t, user)

	//make sure the data type is equal
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.HashPassword, user.HashPassword)

	//make sure auto generated value is not return zero value
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.ChangedPasswordAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.HashPassword, user2.HashPassword)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.ChangedPasswordAt, user2.ChangedPasswordAt, time.Second)
}