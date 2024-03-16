package tokens

import (
	"TrafficPolice/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTokenManager(t *testing.T) {
	signingKey := "sign"

	testCases := []struct {
		name        string
		stringKey   string
		shouldBeErr bool
	}{
		{
			name:        "signingKey is correct. No error",
			stringKey:   signingKey,
			shouldBeErr: false,
		},
		{
			name:        "signingKey is empty. Should be error",
			stringKey:   "",
			shouldBeErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewTokenManager(tc.stringKey)
			assert.Equal(t, tc.shouldBeErr, err != nil)
		})
	}
}

func TestNewJWT(t *testing.T) {
	manager, err := NewTokenManager("sign")
	assert.NoError(t, err)

	testCases := []struct {
		name        string
		tokenInfo   TokenInfo
		ttl         time.Duration
		shouldBeErr bool
	}{
		{
			name:        "Token Info is correct. Should not be err",
			tokenInfo:   TokenInfo{UserID: "some_id", UserRole: domain.ExpertRole},
			ttl:         10 * time.Minute,
			shouldBeErr: false,
		},
		{
			name:        "UserID is empty. Should be err",
			tokenInfo:   TokenInfo{UserRole: domain.ExpertRole},
			ttl:         10 * time.Minute,
			shouldBeErr: true,
		},
		{
			name:        "UserRole is empty. Should be err",
			tokenInfo:   TokenInfo{UserID: "some_id"},
			ttl:         10 * time.Minute,
			shouldBeErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := manager.NewJWT(tc.tokenInfo, tc.ttl)

			assert.Equal(t, tc.shouldBeErr, err != nil)
		})
	}
}
