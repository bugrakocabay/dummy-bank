package db

import (
	"context"
	"database/sql"
	"github.com/bugrakocabay/dummy-bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomEntry(t *testing.T, accountId int64) Entry {
	arg := CreateEntryParams{
		Amount:    util.RandomMoney(),
		AccountID: accountId,
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.Amount, entry.Amount)
	require.Equal(t, arg.AccountID, entry.AccountID)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	acc := createRandomAccount(t)
	createRandomEntry(t, acc.ID)
}

func TestGetEntry(t *testing.T) {
	acc := createRandomAccount(t)
	entry1 := createRandomEntry(t, acc.ID)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	acc := createRandomAccount(t)
	entry1 := createRandomEntry(t, acc.ID)
	arg := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: util.RandomMoney(),
	}
	entry2, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, arg.Amount, entry2.Amount)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	acc := createRandomAccount(t)
	entry1 := createRandomEntry(t, acc.ID)
	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		acc := createRandomAccount(t)
		createRandomEntry(t, acc.ID)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}