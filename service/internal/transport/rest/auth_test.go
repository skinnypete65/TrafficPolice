package rest

import (
	"TrafficPolice/internal/converter"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/service"
	"TrafficPolice/internal/service/mocks"
	"TrafficPolice/internal/transport/rest/dto"
	"bytes"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUp(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	userInfoConverter := converter.NewUserInfoConverter()
	authConverter := converter.NewAuthConverter()
	path := "/auth/sign_up"

	testCases := []struct {
		name             string
		signUp           dto.SignUp
		buildAuthService func() service.AuthService
		expectedCode     int
	}{
		{
			name: "Sign up not existing user. 200 OK",
			signUp: dto.SignUp{
				Username: "user",
				Password: "pass",
			},
			buildAuthService: func() service.AuthService {
				mockService := mocks.NewAuthService(t)
				mockService.On("RegisterExpert", mock.Anything).
					Return(nil)

				return mockService
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Empty username. 400 Bad Request",
			signUp: dto.SignUp{
				Password: "pass",
			},
			buildAuthService: func() service.AuthService {
				mockService := mocks.NewAuthService(t)

				return mockService
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Empty password. 400 Bad Request",
			signUp: dto.SignUp{
				Username: "user",
			},
			buildAuthService: func() service.AuthService {
				mockService := mocks.NewAuthService(t)

				return mockService
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "User with input username already exists. 409 Conflict",
			signUp: dto.SignUp{
				Username: "user",
				Password: "pass",
			},
			buildAuthService: func() service.AuthService {
				mockService := mocks.NewAuthService(t)
				mockService.On("RegisterExpert", mock.Anything).
					Return(errs.ErrAlreadyExists)

				return mockService
			},
			expectedCode: http.StatusConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authService := tc.buildAuthService()
			handler := NewAuthHandler(authService, validate, userInfoConverter, authConverter)

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tc.signUp)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, path, &buf)
			rec := httptest.NewRecorder()

			handler.SignUp(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
