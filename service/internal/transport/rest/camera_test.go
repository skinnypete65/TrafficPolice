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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddCameraType(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	cameraConverter := converter.NewCameraConverter()
	path := "/camera/type"
	cameraTypeID := uuid.New().String()

	testCases := []struct {
		name               string
		input              dto.CameraTypeIn
		buildCameraService func() service.CameraService
		buildAuthService   func() service.AuthService
		expectedCode       int
	}{
		{
			name: "Create camera type. 200 OK",
			input: dto.CameraTypeIn{
				Name: "camerus1",
			},
			buildCameraService: func() service.CameraService {
				mockService := mocks.NewCameraService(t)
				mockService.On("AddCameraType", mock.Anything).
					Return(cameraTypeID, nil)

				return mockService
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Camera type with passed name already exists. 409 Conflict",
			input: dto.CameraTypeIn{
				Name: "camerus1",
			},
			buildCameraService: func() service.CameraService {
				mockService := mocks.NewCameraService(t)
				mockService.On("AddCameraType", mock.Anything).
					Return("", errs.ErrAlreadyExists)

				return mockService
			},
			buildAuthService: func() service.AuthService {
				mockService := mocks.NewAuthService(t)
				return mockService
			},
			expectedCode: http.StatusConflict,
		},
		{
			name:  "Camera type with passed empty name 400 Bad request",
			input: dto.CameraTypeIn{},
			buildCameraService: func() service.CameraService {
				mockService := mocks.NewCameraService(t)

				return mockService
			},
			buildAuthService: func() service.AuthService {
				mockService := mocks.NewAuthService(t)
				return mockService
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewCameraHandler(tc.buildCameraService(), tc.buildAuthService(), validate, cameraConverter)

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tc.input)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, path, &buf)
			rec := httptest.NewRecorder()

			handler.AddCameraType(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestRegisterCamera(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	cameraConverter := converter.NewCameraConverter()
	path := "/camera"
	cameraID := uuid.New().String()
	cameraTypeID := uuid.New().String()

	testCases := []struct {
		name               string
		input              dto.RegisterCamera
		buildCameraService func() service.CameraService
		buildAuthService   func() service.AuthService
		expectedCode       int
	}{
		{
			name: "Create camera. 200 OK",
			input: dto.RegisterCamera{
				CameraIn: dto.CameraIn{
					CameraTypeID: cameraTypeID,
					Latitude:     30.0,
					Longitude:    30.0,
					ShortDesc:    "short desc",
				},
				SignUp: dto.SignUp{
					Username: "user",
					Password: "pass",
				},
			},
			buildCameraService: func() service.CameraService {
				mockService := mocks.NewCameraService(t)
				return mockService
			},
			buildAuthService: func() service.AuthService {
				mockService := mocks.NewAuthService(t)
				mockService.On("RegisterCamera", mock.Anything).
					Return(cameraID, nil)

				return mockService
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Camera with input username already exists. 409 Conflict",
			input: dto.RegisterCamera{
				CameraIn: dto.CameraIn{
					CameraTypeID: cameraTypeID,
					Latitude:     30.0,
					Longitude:    30.0,
					ShortDesc:    "short desc",
				},
				SignUp: dto.SignUp{
					Username: "user",
					Password: "pass",
				},
			},
			buildCameraService: func() service.CameraService {
				mockService := mocks.NewCameraService(t)
				return mockService
			},
			buildAuthService: func() service.AuthService {
				mockService := mocks.NewAuthService(t)
				mockService.On("RegisterCamera", mock.Anything).
					Return(cameraID, nil)

				return mockService
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewCameraHandler(tc.buildCameraService(), tc.buildAuthService(), validate, cameraConverter)

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tc.input)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, path, &buf)
			rec := httptest.NewRecorder()

			handler.RegisterCamera(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
