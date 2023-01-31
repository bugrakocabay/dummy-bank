package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	mockdb "github.com/bugrakocabay/dummy-bank/db/mock"
	db "github.com/bugrakocabay/dummy-bank/db/sqlc"
	"github.com/bugrakocabay/dummy-bank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTransferAPI(t *testing.T) {
	account1 := createRandomAccount()
	account2 := createRandomAccount()
	transfer := createRandomTransfer(account1.ID, account2.ID)

	testCases := []struct {
		name          string
		transferID    int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			transferID: transfer.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTransfer(gomock.Any(), gomock.Eq(transfer.ID)).
					Times(1).
					Return(transfer, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTransfer(t, recorder.Body, transfer)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/transfers/%d", tc.transferID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			fmt.Println(recorder)
			tc.checkResponse(t, recorder)
		})

	}

}

func createRandomTransfer(account1ID, account2ID int64) db.Transfer {
	return db.Transfer{
		ID:            util.RandomInt(0, 1000),
		FromAccountID: account1ID,
		ToAccountID:   account2ID,
		Amount:        2,
	}
}

func requireBodyMatchTransfer(t *testing.T, body *bytes.Buffer, transfer db.Transfer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTransfer db.Transfer
	err = json.Unmarshal(data, &gotTransfer)
	require.NoError(t, err)
	require.NotEqual(t, transfer, gotTransfer)
}
