package v1

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"refactoring/internal/api/v1/mocks"
	"refactoring/internal/model"
	storageErr "refactoring/internal/storage/errors"
)

func TestAPI_searchUsers(t *testing.T) {

	tests := []struct {
		name         string
		mockStorage  *mocks.Storage
		expectedCode int
	}{
		{
			name: "OK",
			mockStorage: func() *mocks.Storage {
				mockStorage := mocks.Storage{}
				mockStorage.On("GetAllUsers").
					Return(&model.UserStore{
						Increment: 1,
						List:      model.UserList{"1": {CreatedAt: time.Now(), DisplayName: "testName", Email: "testEmail"}}}, nil)
				return &mockStorage
			}(),
			expectedCode: http.StatusOK,
		},
		{
			name: "unexpected err on getting users",
			mockStorage: func() *mocks.Storage {
				mockStorage := mocks.Storage{}
				mockStorage.On("GetAllUsers").
					Return(nil, errors.New("something err"))
				return &mockStorage
			}(),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "no content",
			mockStorage: func() *mocks.Storage {
				mockStorage := mocks.Storage{}
				mockStorage.On("GetAllUsers").
					Return(&model.UserStore{
						Increment: 0,
						List:      model.UserList{}}, nil)
				return &mockStorage
			}(),
			expectedCode: http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAPI := API{}
			testAPI.storage = tt.mockStorage

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()

			testAPI.searchUsers(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestAPI_getUser(t *testing.T) {

	tests := []struct {
		name         string
		mockStorage  *mocks.Storage
		id           string
		expectedCode int
	}{
		{
			name: "OK",
			mockStorage: func() *mocks.Storage {
				mockStorage := mocks.Storage{}
				mockStorage.On("GetUser", int64(1)).
					Return(&model.User{
						CreatedAt:   time.Unix(1, 1),
						DisplayName: "testName",
						Email:       "testEmail",
					}, nil)
				return &mockStorage
			}(),
			id:           "1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "bad id",
			id:           "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "user is not found",
			mockStorage: func() *mocks.Storage {
				mockStorage := mocks.Storage{}
				mockStorage.On("GetUser", int64(1)).
					Return(nil, storageErr.ErrUserIsNotFound)
				return &mockStorage
			}(),
			id:           "1",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "unexpected err",
			mockStorage: func() *mocks.Storage {
				mockStorage := mocks.Storage{}
				mockStorage.On("GetUser", int64(1)).
					Return(nil, errors.New("unexpected error"))
				return &mockStorage
			}(),
			id:           "1",
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAPI := API{}
			testAPI.storage = tt.mockStorage

			routCtx := chi.NewRouteContext()
			routCtx.URLParams.Add("id", tt.id)

			req := httptest.NewRequest(http.MethodGet, "/test", nil).
				WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, routCtx))
			rec := httptest.NewRecorder()

			testAPI.getUser(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}
