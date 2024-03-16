package middlewares

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/service/mocks"
	"TrafficPolice/internal/tokens"
	"bytes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
