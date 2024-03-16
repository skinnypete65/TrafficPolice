package middlewares

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/service"
	"TrafficPolice/internal/service/mocks"
	"TrafficPolice/internal/tokens"
	"bytes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestIdentifyRole(t *testing.T) {
	path := "/auth"
	signingKey := "sign"
	tokenManager, err := tokens.NewTokenManager(signingKey)
	assert.NoError(t, err)

	expertService := mocks.NewExpertService(t)

	userID := uuid.New().String()
	testCases := []struct {
		name         string
		role         domain.Role
		buildHeader  func() string
		headerKey    string
		expectedCode int
	}{
		{
			name: "Pass correct token and role. 200 OK",
			role: domain.DirectorRole,
			buildHeader: func() string {
				token, err := tokenManager.NewJWT(
					tokens.TokenInfo{UserID: userID, UserRole: domain.DirectorRole},
					10*time.Hour,
				)
				assert.NoError(t, err)
				return "Bearer " + token
			},
			headerKey:    authorizationHeader,
			expectedCode: http.StatusOK,
		},
		{
			name: "Pass correct token but not correct role. 401 Unauthorized",
			role: domain.DirectorRole,
			buildHeader: func() string {
				token, err := tokenManager.NewJWT(
					tokens.TokenInfo{UserID: userID, UserRole: domain.ExpertRole},
					10*time.Hour,
				)
				assert.NoError(t, err)
				return "Bearer " + token
			},
			headerKey:    authorizationHeader,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Pass incorrect header. 401 Unauthorized",
			role: domain.DirectorRole,
			buildHeader: func() string {
				token, err := tokenManager.NewJWT(
					tokens.TokenInfo{UserID: userID, UserRole: domain.ExpertRole},
					10*time.Hour,
				)
				assert.NoError(t, err)
				return token
			},
			headerKey:    authorizationHeader,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Pass empty header. 401 Unauthorized",
			role: domain.DirectorRole,
			buildHeader: func() string {
				return ""
			},
			headerKey:    authorizationHeader,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Pass empty token. 401 Unauthorized",
			role: domain.DirectorRole,
			buildHeader: func() string {
				return "Bearer "
			},
			headerKey:    authorizationHeader,
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			middleware := NewAuthMiddleware(tokenManager, expertService)
			var buf bytes.Buffer

			req := httptest.NewRequest(http.MethodGet, path, &buf)
			rec := httptest.NewRecorder()

			header := tc.buildHeader()
			req.Header.Set(tc.headerKey, header)

			middleware.IdentifyRole(http.HandlerFunc(okHandler), tc.role).
				ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestIsExpertConfirmed(t *testing.T) {
	path := "/auth"
	signingKey := "sign"
	tokenManager, err := tokens.NewTokenManager(signingKey)
	assert.NoError(t, err)

	userID := uuid.New().String()
	testCases := []struct {
		name               string
		buildHeader        func() string
		headerKey          string
		buildExpertService func() service.ExpertService
		expectedCode       int
	}{
		{
			name: "Pass correct token and role. Expert is confirmed. 200 OK",
			buildHeader: func() string {
				token, err := tokenManager.NewJWT(
					tokens.TokenInfo{UserID: userID, UserRole: domain.ExpertRole},
					10*time.Hour,
				)
				assert.NoError(t, err)
				return "Bearer " + token
			},
			headerKey: authorizationHeader,
			buildExpertService: func() service.ExpertService {
				expertService := mocks.NewExpertService(t)
				expertService.On("GetExpertByUserID", mock.Anything).
					Return(domain.Expert{IsConfirmed: true}, nil)

				return expertService
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Pass incorrect header. 401 Unauthorized",
			buildHeader: func() string {
				token, err := tokenManager.NewJWT(
					tokens.TokenInfo{UserID: userID, UserRole: domain.ExpertRole},
					10*time.Hour,
				)
				assert.NoError(t, err)
				return token
			},
			headerKey: authorizationHeader,
			buildExpertService: func() service.ExpertService {
				expertService := mocks.NewExpertService(t)
				return expertService
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Pass empty header. 401 Unauthorized",
			buildHeader: func() string {
				return ""
			},
			headerKey: authorizationHeader,
			buildExpertService: func() service.ExpertService {
				expertService := mocks.NewExpertService(t)
				return expertService
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Pass correct token and not expert role. 200 OK",
			buildHeader: func() string {
				token, err := tokenManager.NewJWT(
					tokens.TokenInfo{UserID: userID, UserRole: domain.DirectorRole},
					10*time.Hour,
				)
				assert.NoError(t, err)
				return "Bearer " + token
			},
			headerKey: authorizationHeader,
			buildExpertService: func() service.ExpertService {
				expertService := mocks.NewExpertService(t)

				return expertService
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Pass correct token and role. Expert is not confirmed. 403 Forbidden",
			buildHeader: func() string {
				token, err := tokenManager.NewJWT(
					tokens.TokenInfo{UserID: userID, UserRole: domain.ExpertRole},
					10*time.Hour,
				)
				assert.NoError(t, err)
				return "Bearer " + token
			},
			headerKey: authorizationHeader,
			buildExpertService: func() service.ExpertService {
				expertService := mocks.NewExpertService(t)
				expertService.On("GetExpertByUserID", mock.Anything).
					Return(domain.Expert{IsConfirmed: false}, nil)

				return expertService
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "Expert with user id not exists. 401 Unauthorized",
			buildHeader: func() string {
				token, err := tokenManager.NewJWT(
					tokens.TokenInfo{UserID: "wrong_user_id", UserRole: domain.ExpertRole},
					10*time.Hour,
				)
				assert.NoError(t, err)
				return "Bearer " + token
			},
			headerKey: authorizationHeader,
			buildExpertService: func() service.ExpertService {
				expertService := mocks.NewExpertService(t)
				expertService.On("GetExpertByUserID", mock.Anything).
					Return(domain.Expert{}, errs.ErrUserNotExists)

				return expertService
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			middleware := NewAuthMiddleware(tokenManager, tc.buildExpertService())
			var buf bytes.Buffer

			req := httptest.NewRequest(http.MethodGet, path, &buf)
			rec := httptest.NewRecorder()

			header := tc.buildHeader()
			req.Header.Set(tc.headerKey, header)

			middleware.IsExpertConfirmed(http.HandlerFunc(okHandler)).
				ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
