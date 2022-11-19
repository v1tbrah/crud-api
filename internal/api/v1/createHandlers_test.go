package v1

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"refactoring/internal/api/v1/mocks"
)

func TestAPI_createUser(t *testing.T) {
	tests := []struct {
		name         string
		payload      string
		mockStorage  *mocks.Storage
		id           string
		expectedCode int
	}{
		{
			name:    "OK",
			payload: "{\"display_name\": \"testName\", \"email\": \"testEmail\"}",
			mockStorage: func() *mocks.Storage {
				mockStorage := mocks.Storage{}
				mockStorage.On("CreateUser", mock.AnythingOfType("*model.User")).
					Return(int64(1), nil)
				return &mockStorage
			}(),
			expectedCode: http.StatusCreated,
		},
		{
			name:         "bad payload",
			payload:      "{\"display_name\": ___",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "unexpected err on creating user",
			payload: "{\"display_name\": \"testName\", \"email\": \"testEmail\"}",
			mockStorage: func() *mocks.Storage {
				mockStorage := mocks.Storage{}
				mockStorage.On("CreateUser", mock.AnythingOfType("*model.User")).
					Return(int64(0), errors.New("unexpected error"))
				return &mockStorage
			}(),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAPI := API{}
			testAPI.storage = tt.mockStorage

			b := &bytes.Buffer{}
			b.WriteString(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/test", b)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			testAPI.createUser(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}
