package service

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/repository"
	"TrafficPolice/internal/repository/mocks"
	"TrafficPolice/internal/tokens"
	"TrafficPolice/pkg/hash"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func newHasher(salt string) *hash.SHA1Hasher {
	return hash.NewSHA1Hasher(salt)
}

func TestRegisterExpert(t *testing.T) {
	hasher := newHasher("salt")

	testCases := []struct {
		name            string
		buildAuthRepo   func() repository.AuthRepo
		buildRatingRepo func() repository.RatingRepo
		buildUserInfo   func() domain.UserInfo
		expectedErr     error
	}{
		{
			name: "Should be registered expert, no error",
			buildAuthRepo: func() repository.AuthRepo {
				mockRepo := mocks.NewAuthRepo(t)

				mockRepo.On("CheckUserExists", mock.Anything).
					Return(false)
				mockRepo.On("InsertUser", mock.Anything).
					Return(nil)
				mockRepo.On("InsertExpert", mock.Anything).
					Return(nil)

				return mockRepo
			},
			buildRatingRepo: func() repository.RatingRepo {
				return mocks.NewRatingRepo(t)
			},
			buildUserInfo: func() domain.UserInfo {
				return domain.UserInfo{}
			},
			expectedErr: nil,
		},
		{
			name: "User already exists, expected ErrAlreadyExists error",
			buildAuthRepo: func() repository.AuthRepo {
				mockRepo := mocks.NewAuthRepo(t)
				mockRepo.On("CheckUserExists", mock.Anything).
					Return(true)

				return mockRepo
			},
			buildRatingRepo: func() repository.RatingRepo {
				return mocks.NewRatingRepo(t)
			},
			buildUserInfo: func() domain.UserInfo {
				return domain.UserInfo{}
			},
			expectedErr: errs.ErrAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authRepo := tc.buildAuthRepo()
			ratingRepo := tc.buildRatingRepo()

			tokenManager, err := tokens.NewTokenManager("sign")
			assert.NoError(t, err)

			authService := NewAuthService(authRepo, ratingRepo, hasher, tokenManager)

			err = authService.RegisterExpert(tc.buildUserInfo())
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRegisterCamera(t *testing.T) {
	hasher := newHasher("salt")

	testCases := []struct {
		name                string
		buildAuthRepo       func() repository.AuthRepo
		buildRatingRepo     func() repository.RatingRepo
		buildRegisterCamera func() domain.RegisterCamera
		expectedErr         error
	}{
		{
			name: "Should be registered camera, no error",
			buildAuthRepo: func() repository.AuthRepo {
				mockRepo := mocks.NewAuthRepo(t)

				mockRepo.On("CheckUserExists", mock.Anything).
					Return(false)
				mockRepo.On("InsertUser", mock.Anything).
					Return(nil)
				mockRepo.On("InsertCamera", mock.Anything, mock.Anything).
					Return("some id", nil)

				return mockRepo
			},
			buildRatingRepo: func() repository.RatingRepo {
				return mocks.NewRatingRepo(t)
			},
			buildRegisterCamera: func() domain.RegisterCamera {
				return domain.RegisterCamera{}
			},
			expectedErr: nil,
		},
		{
			name: "Camera with username already exists. Expect ErrAlreadyExists",
			buildAuthRepo: func() repository.AuthRepo {
				mockRepo := mocks.NewAuthRepo(t)

				mockRepo.On("CheckUserExists", mock.Anything).
					Return(true)

				return mockRepo
			},
			buildRatingRepo: func() repository.RatingRepo {
				return mocks.NewRatingRepo(t)
			},
			buildRegisterCamera: func() domain.RegisterCamera {
				return domain.RegisterCamera{}
			},
			expectedErr: errs.ErrAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authRepo := tc.buildAuthRepo()
			ratingRepo := tc.buildRatingRepo()

			tokenManager, err := tokens.NewTokenManager("sign")
			assert.NoError(t, err)

			authService := NewAuthService(authRepo, ratingRepo, hasher, tokenManager)

			_, err = authService.RegisterCamera(tc.buildRegisterCamera())
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRegisterDirectors(t *testing.T) {
	hasher := newHasher("salt")

	testCases := []struct {
		name            string
		buildAuthRepo   func() repository.AuthRepo
		buildRatingRepo func() repository.RatingRepo
		buildUsers      func() []domain.UserInfo
		expectedErr     error
	}{
		{
			name: "Should be registered all directors",
			buildAuthRepo: func() repository.AuthRepo {
				mockRepo := mocks.NewAuthRepo(t)

				mockRepo.On("CheckUserExists", mock.Anything).
					Return(false)
				mockRepo.On("InsertUser", mock.Anything).
					Return(nil)
				mockRepo.On("InsertDirector", mock.Anything, mock.Anything).
					Return(nil)

				return mockRepo
			},
			buildRatingRepo: func() repository.RatingRepo {
				return mocks.NewRatingRepo(t)
			},
			buildUsers: func() []domain.UserInfo {
				return []domain.UserInfo{
					{Username: "director1", Password: "director1"},
					{Username: "director2", Password: "director2"},
				}
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authRepo := tc.buildAuthRepo()
			ratingRepo := tc.buildRatingRepo()

			tokenManager, err := tokens.NewTokenManager("sign")
			assert.NoError(t, err)

			authService := NewAuthService(authRepo, ratingRepo, hasher, tokenManager)

			err = authService.RegisterDirectors(tc.buildUsers())
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSignIn(t *testing.T) {
	hasher := newHasher("salt")
	tokenManager, err := tokens.NewTokenManager("sign")
	assert.NoError(t, err)

	userUUID := uuid.New()
	buildInputUser := func(pass string) domain.UserInfo {
		return domain.UserInfo{
			ID:       userUUID,
			Username: "user",
			Password: pass,
			UserRole: string(domain.ExpertRole),
		}
	}

	buildOutputUser := func(pass string) domain.UserInfo {
		hashedPass, err := hasher.Hash(pass)
		assert.NoError(t, err)

		return domain.UserInfo{
			ID:       userUUID,
			Username: "user",
			Password: hashedPass,
			UserRole: string(domain.ExpertRole),
		}
	}

	testCases := []struct {
		name            string
		buildAuthRepo   func() repository.AuthRepo
		buildRatingRepo func() repository.RatingRepo
		buildUser       func() domain.UserInfo
		expectedErr     error
	}{
		{
			name: "User should be signed up. No error",
			buildAuthRepo: func() repository.AuthRepo {
				mockRepo := mocks.NewAuthRepo(t)

				mockRepo.On("SignIn", mock.Anything).
					Return(buildOutputUser("user_password"), nil)

				return mockRepo
			},
			buildRatingRepo: func() repository.RatingRepo {
				return mocks.NewRatingRepo(t)
			},
			buildUser: func() domain.UserInfo {
				return buildInputUser("user_password")
			},
			expectedErr: nil,
		},
		{
			name: "Invalid username. Expected ErrNoRows",
			buildAuthRepo: func() repository.AuthRepo {
				mockRepo := mocks.NewAuthRepo(t)

				mockRepo.On("SignIn", mock.Anything).
					Return(domain.UserInfo{}, errs.ErrNoRows)

				return mockRepo
			},
			buildRatingRepo: func() repository.RatingRepo {
				return mocks.NewRatingRepo(t)
			},
			buildUser: func() domain.UserInfo {
				return buildInputUser("user_password")
			},
			expectedErr: errs.ErrNoRows,
		},
		{
			name: "Invalid password. Expected ErrInvalidPass",
			buildAuthRepo: func() repository.AuthRepo {
				mockRepo := mocks.NewAuthRepo(t)

				mockRepo.On("SignIn", mock.Anything).
					Return(buildOutputUser("right_password"), nil)

				return mockRepo
			},
			buildRatingRepo: func() repository.RatingRepo {
				return mocks.NewRatingRepo(t)
			},
			buildUser: func() domain.UserInfo {
				return buildInputUser("wrong_password")
			},
			expectedErr: errs.ErrInvalidPass,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authRepo := tc.buildAuthRepo()
			ratingRepo := tc.buildRatingRepo()

			authService := NewAuthService(authRepo, ratingRepo, hasher, tokenManager)

			_, err = authService.SignIn(tc.buildUser())
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestConfirmExpert(t *testing.T) {
	hasher := newHasher("salt")
	errConfirmExpert := errors.New("error while confirming expert")

	testCases := []struct {
		name               string
		buildAuthRepo      func() repository.AuthRepo
		buildRatingRepo    func() repository.RatingRepo
		buildConfirmExpert func() domain.ConfirmExpert
		expectedErr        error
	}{
		{
			name: "Expert should be confirmed. Expect no error",
			buildAuthRepo: func() repository.AuthRepo {
				mockRepo := mocks.NewAuthRepo(t)

				mockRepo.On("ConfirmExpert", mock.Anything).
					Return(nil)

				return mockRepo
			},
			buildRatingRepo: func() repository.RatingRepo {
				mockRepo := mocks.NewRatingRepo(t)
				mockRepo.On("InsertExpertId", mock.Anything).
					Return(nil)

				return mockRepo
			},
			buildConfirmExpert: func() domain.ConfirmExpert {
				return domain.ConfirmExpert{
					ExpertID:    uuid.New().String(),
					IsConfirmed: true,
				}
			},
			expectedErr: nil,
		},
		{
			name: "Problems while confirm expert. Expect error",
			buildAuthRepo: func() repository.AuthRepo {
				mockRepo := mocks.NewAuthRepo(t)

				mockRepo.On("ConfirmExpert", mock.Anything).
					Return(errConfirmExpert)

				return mockRepo
			},
			buildRatingRepo: func() repository.RatingRepo {
				mockRepo := mocks.NewRatingRepo(t)
				return mockRepo
			},
			buildConfirmExpert: func() domain.ConfirmExpert {
				return domain.ConfirmExpert{
					ExpertID:    uuid.New().String(),
					IsConfirmed: true,
				}
			},
			expectedErr: errConfirmExpert,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authRepo := tc.buildAuthRepo()
			ratingRepo := tc.buildRatingRepo()

			tokenManager, err := tokens.NewTokenManager("sign")
			assert.NoError(t, err)

			authService := NewAuthService(authRepo, ratingRepo, hasher, tokenManager)

			err = authService.ConfirmExpert(tc.buildConfirmExpert())
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
