package rest

import (
	"TrafficPolice/internal/converter"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/service"
	"TrafficPolice/internal/service/mocks"
	"TrafficPolice/internal/transport/rest/dto"
	"bytes"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

func TestSignIn(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	userInfoConverter := converter.NewUserInfoConverter()
	authConverter := converter.NewAuthConverter()
	path := "/auth/sign_in"

	testCases := []struct {
		name             string
		signIn           dto.SignInInput
		buildAuthService func() service.AuthService
		expectedCode     int
	}{
		{
			name: "Sign in with correct username and password. 200 OK",
			signIn: dto.SignInInput{
				Username: "user",
				Password: "pass",
			},
			buildAuthService: func() service.AuthService {
				mockService := mocks.NewAuthService(t)
				mockService.On("SignIn", mock.Anything).
					Return(domain.Tokens{AccessToken: "access", RefreshToken: "refresh"}, nil)

				return mockService
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Sign in with blank username. 400 Bad request",
			signIn: dto.SignInInput{
				Password: "pass",
			},
			buildAuthService: func() service.AuthService {
				mockService := mocks.NewAuthService(t)

				return mockService
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Sign in with blank password. 400 Bad request",
			signIn: dto.SignInInput{
				Username: "user",
			},
			buildAuthService: func() service.AuthService {
				mockService := mocks.NewAuthService(t)

				return mockService
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Sign in with not existing user. 404 Not found",
			signIn: dto.SignInInput{
				Username: "user",
				Password: "pass",
			},
			buildAuthService: func() service.AuthService {
				mockService := mocks.NewAuthService(t)
				mockService.On("SignIn", mock.Anything).
					Return(domain.Tokens{}, errs.ErrNoRows)

				return mockService
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "Sign in with wrong password. 401 Unauthorized",
			signIn: dto.SignInInput{
				Username: "user",
				Password: "pass",
			},
			buildAuthService: func() service.AuthService {
				mockService := mocks.NewAuthService(t)
				mockService.On("SignIn", mock.Anything).
					Return(domain.Tokens{}, errs.ErrInvalidPass)

				return mockService
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authService := tc.buildAuthService()
			handler := NewAuthHandler(authService, validate, userInfoConverter, authConverter)

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tc.signIn)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, path, &buf)
			rec := httptest.NewRecorder()

			handler.SignIn(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestConfirmExpert(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	userInfoConverter := converter.NewUserInfoConverter()
	authConverter := converter.NewAuthConverter()
	path := "/auth/confirm/expert"
	expertID := uuid.New().String()

	testCases := []struct {
		name             string
		input            dto.ConfirmExpertInput
		buildAuthService func() service.AuthService
		expectedCode     int
	}{
		{
			name: "Confirm expert. 200 OK",
			input: dto.ConfirmExpertInput{
				ExpertID:    expertID,
				IsConfirmed: true,
			},
			buildAuthService: func() service.AuthService {
				mockService := mocks.NewAuthService(t)
				mockService.On("ConfirmExpert", mock.Anything).
					Return(nil)

				return mockService
			},
			expectedCode: http.StatusOK,
		},
		{
			name:  "Expert id not passed. 400 Bad request",
			input: dto.ConfirmExpertInput{},
			buildAuthService: func() service.AuthService {
				mockService := mocks.NewAuthService(t)

				return mockService
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Expert with passed id not exists. 404 Not found",
			input: dto.ConfirmExpertInput{
				ExpertID:    expertID,
				IsConfirmed: true,
			},
			buildAuthService: func() service.AuthService {
				mockService := mocks.NewAuthService(t)
				mockService.On("ConfirmExpert", mock.Anything).
					Return(errs.ErrNoRows)

				return mockService
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authService := tc.buildAuthService()
			handler := NewAuthHandler(authService, validate, userInfoConverter, authConverter)

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tc.input)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, path, &buf)
			rec := httptest.NewRecorder()

			handler.ConfirmExpert(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
