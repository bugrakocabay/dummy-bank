package db

import (
	"context"
	"github.com/bugrakocabay/dummy-bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	arg := CreateTransferParams{
		Amount:        util.RandomMoney(),
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.Amount, transfer.Amount)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	createRandomTransfer(t, acc1, acc2)
}

func TestGetTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, acc1, acc2)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		acc1 := createRandomAccount(t)
		acc2 := createRandomAccount(t)
		createRandomTransfer(t, acc1, acc2)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}