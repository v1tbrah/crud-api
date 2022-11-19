package v1

import (
	"bytes"
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

func TestAPI_updateUser(t *testing.T) {
	tests := []struct {
		name         string
		payload      string
		mockStorage  *mocks.Storage
		id           string
		expectedCode int
	}{
		{
			name:    "OK",
			payload: `{"display_name": "testName"}`,
			mockStorage: func() *mocks.Storage {
				mockStorage := mocks.Storage{}
				mockStorage.On("UpdateUser", int64(1), "testName").
					Return(nil)
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
			name:         "bad payload",
			payload:      "{\"display_name\": ___",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "user is not found",
			payload: `{"display_name": "testName"}`,
			mockStorage: func() *mocks.Storage {
				mockStorage := mocks.Storage{}
				mockStorage.On("UpdateUser", int64(1), "testName").
					Return(storageErr.ErrUserIsNotFound)
				return &mockStorage
			}(),
			id:           "1",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "unexpected error on updating user",
			payload: `{"display_name": "testName"}`,
			mockStorage: func() *mocks.Storage {
				mockStorage := mocks.Storage{}
				mockStorage.On("UpdateUser", int64(1), "testName").
					Return(errors.New("unexpected error"))
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

			b := &bytes.Buffer{}
			b.WriteString(tt.payload)
			req := httptest.NewRequest(http.MethodPatch, "/test", b).
				WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, routCtx))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			testAPI.updateUser(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}
