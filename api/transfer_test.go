package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/bugrakocabay/dummy-bank/db/mock"
	db "github.com/bugrakocabay/dummy-bank/db/sqlc"
	"github.com/bugrakocabay/dummy-bank/token"
	"github.com/bugrakocabay/dummy-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetTransferAPI(t *testing.T) {
	user1, _ := createRandomUser(t)
	user2, _ := createRandomUser(t)

	account1 := createRandomAccount(user1.Username)
	account2 := createRandomAccount(user2.Username)
	transfer := createRandomTransfer(account1.ID, account2.ID)

	testCases := []struct {
		name          string
		transferID    int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			transferID: transfer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTransfer(gomock.Any(), gomock.Eq(transfer.ID)).
					Times(1).
					Return(transfer, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:       "NotFound",
			transferID: transfer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTransfer(gomock.Any(), gomock.Eq(transfer.ID)).
					Times(1).
					Return(db.Transfer{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalServerErr",
			transferID: transfer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTransfer(gomock.Any(), gomock.Eq(transfer.ID)).
					Times(1).
					Return(db.Transfer{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "BadRequest",
			transferID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTransfer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/transfers/%d", tc.transferID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func TestCreateTransferAPI(t *testing.T) {
	user1, _ := createRandomUser(t)
	user2, _ := createRandomUser(t)

	account1 := createRandomAccount(user1.Username)
	account2 := createRandomAccount(user2.Username)
	account1.Currency = "USD"
	account2.Currency = "USD"

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Created",
			body: gin.H{"from_account_id": account1.ID, "to_account_id": account2.ID, "amount": int64(1), "currency": "USD"},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        1,
				}
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

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
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/transfers")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
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
