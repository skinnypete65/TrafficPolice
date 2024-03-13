package services

import (
	"TrafficPolice/errs"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"TrafficPolice/internal/tokens"
	"TrafficPolice/pkg/hash"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type AuthService interface {
	RegisterExpert(input domain.UserInfo) error
	RegisterCamera(info domain.RegisterCamera) error
	RegisterDirectors(users []domain.UserInfo) error
	SignIn(input domain.UserInfo) (string, error)
	ConfirmExpert(data domain.ConfirmExpert) error
	ParseAccessToken(accessToken string) (tokens.TokenInfo, error)
}

type authService struct {
	repo           repository.AuthRepo
	hasher         hash.PasswordHasher
	tokenManager   tokens.TokenManager
	accessTokenTTL time.Duration
}

func NewAuthService(repo repository.AuthRepo, tokenManager tokens.TokenManager, passSalt string) AuthService {
	return &authService{
		repo:           repo,
		hasher:         hash.NewSHA1Hasher(passSalt),
		tokenManager:   tokenManager,
		accessTokenTTL: 30 * 24 * time.Hour,
	}
}

func (s *authService) RegisterExpert(user domain.UserInfo) error {
	alreadyExists := s.repo.CheckUserExists(user.Username)
	if alreadyExists {
		return errs.ErrAlreadyExists
	}

	hashedPass, err := s.hasher.Hash(user.Password)
	if err != nil {
		return err
	}

	user.ID = uuid.New()
	user.Password = hashedPass
	user.UserRole = "expert"
	err = s.repo.InsertUser(user)
	if err != nil {
		return err
	}

	err = s.repo.InsertExpert(domain.Expert{
		ID:       uuid.New().String(),
		UserInfo: user,
	})

	return err
}

func (s *authService) RegisterCamera(info domain.RegisterCamera) error {
	alreadyExists := s.repo.CheckUserExists(info.Username)
	if alreadyExists {
		return fmt.Errorf("camera with username '%s' already exists", info.Username)
	}

	hashedPass, err := s.hasher.Hash(info.Password)
	if err != nil {
		return err
	}

	userID := uuid.New()
	info.Camera.ID = uuid.New().String()

	userInfo := domain.UserInfo{
		ID:       userID,
		Username: info.Username,
		Password: hashedPass,
		UserRole: string(domain.CameraRole),
	}
	err = s.repo.InsertUser(userInfo)
	if err != nil {
		return err
	}

	return s.repo.InsertCamera(info.Camera, userID)
}

func (s *authService) RegisterDirectors(users []domain.UserInfo) error {
	for i := range users {
		alreadyExists := s.repo.CheckUserExists(users[i].Username)
		if alreadyExists {
			continue
		}

		hashedPass, err := s.hasher.Hash(users[i].Password)
		if err != nil {
			return err
		}

		users[i].ID = uuid.New()
		users[i].Password = hashedPass
		users[i].UserRole = "director"

		err = s.repo.InsertUser(users[i])
		if err != nil {
			return err
		}

		err = s.repo.InsertDirector(domain.Director{
			ID:   uuid.New(),
			User: users[i],
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *authService) SignIn(input domain.UserInfo) (string, error) {
	user, err := s.repo.SignIn(input.Username)
	if err != nil {
		return "", err
	}

	inputHashPass, err := s.hasher.Hash(input.Password)
	if err != nil {
		return "", err
	}

	if user.Password != inputHashPass {
		return "", errs.ErrInvalidPass
	}

	return s.tokenManager.NewJWT(tokens.TokenInfo{
		UserID:   user.ID.String(),
		UserRole: domain.Role(user.UserRole),
	}, s.accessTokenTTL)
}

func (s *authService) ConfirmExpert(data domain.ConfirmExpert) error {
	return s.repo.ConfirmExpert(data)
}

func (s *authService) ParseAccessToken(accessToken string) (tokens.TokenInfo, error) {
	return s.tokenManager.Parse(accessToken)
}
