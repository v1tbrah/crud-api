package v1

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"refactoring/internal/api/v1/mocks"
	storageErr "refactoring/internal/storage/errors"
)

func TestAPI_deleteUser(t *testing.T) {
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
				mockStorage.On("DeleteUser", int64(1)).Return(nil)
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
				mockStorage.On("DeleteUser", int64(1)).Return(storageErr.ErrUserIsNotFound)
				return &mockStorage
			}(),
			id:           "1",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "unexpected err",
			mockStorage: func() *mocks.Storage {
				mockStorage := mocks.Storage{}
				mockStorage.On("DeleteUser", int64(1)).Return(errors.New("unexpected error"))
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

			req := httptest.NewRequest(http.MethodDelete, "/test", nil).
				WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, routCtx))
			rec := httptest.NewRecorder()

			testAPI.deleteUser(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}
