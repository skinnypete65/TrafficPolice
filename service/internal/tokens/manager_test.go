package tokens

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
