package services

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/repository"
	"TrafficPolice/internal/tokens"
	"TrafficPolice/pkg/hash"
	"github.com/google/uuid"
	"time"
)

type AuthService interface {
	RegisterExpert(input domain.UserInfo) error
	RegisterCamera(info domain.RegisterCamera) (string, error)
	RegisterDirectors(users []domain.UserInfo) error
	SignIn(input domain.UserInfo) (domain.Tokens, error)
	ConfirmExpert(data domain.ConfirmExpert) error
	ParseAccessToken(accessToken string) (tokens.TokenInfo, error)
}

type authService struct {
	authRepo       repository.AuthRepo
	ratingRepo     repository.RatingRepo
	hasher         hash.PasswordHasher
	tokenManager   tokens.TokenManager
	accessTokenTTL time.Duration
}

func NewAuthService(
	repo repository.AuthRepo,
	ratingRepo repository.RatingRepo,
	tokenManager tokens.TokenManager,
	passSalt string,
) AuthService {
	return &authService{
		authRepo:       repo,
		ratingRepo:     ratingRepo,
		hasher:         hash.NewSHA1Hasher(passSalt),
		tokenManager:   tokenManager,
		accessTokenTTL: 30 * 24 * time.Hour,
	}
}

func (s *authService) RegisterExpert(user domain.UserInfo) error {
	alreadyExists := s.authRepo.CheckUserExists(user.Username)
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
	err = s.authRepo.InsertUser(user)
	if err != nil {
		return err
	}

	err = s.authRepo.InsertExpert(domain.Expert{
		ID:       uuid.New().String(),
		UserInfo: user,
	})

	return err
}

func (s *authService) RegisterCamera(info domain.RegisterCamera) (string, error) {
	alreadyExists := s.authRepo.CheckUserExists(info.Username)
	if alreadyExists {
		return "", errs.ErrAlreadyExists
	}

	hashedPass, err := s.hasher.Hash(info.Password)
	if err != nil {
		return "", err
	}

	userID := uuid.New()
	info.Camera.ID = uuid.New().String()

	userInfo := domain.UserInfo{
		ID:       userID,
		Username: info.Username,
		Password: hashedPass,
		UserRole: string(domain.CameraRole),
	}
	err = s.authRepo.InsertUser(userInfo)
	if err != nil {
		return "", err
	}

	return s.authRepo.InsertCamera(info.Camera, userID)
}

func (s *authService) RegisterDirectors(users []domain.UserInfo) error {
	for i := range users {
		alreadyExists := s.authRepo.CheckUserExists(users[i].Username)
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

		err = s.authRepo.InsertUser(users[i])
		if err != nil {
			return err
		}

		err = s.authRepo.InsertDirector(domain.Director{
			ID:   uuid.New(),
			User: users[i],
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *authService) SignIn(input domain.UserInfo) (domain.Tokens, error) {
	user, err := s.authRepo.SignIn(input.Username)
	if err != nil {
		return domain.Tokens{}, err
	}

	inputHashPass, err := s.hasher.Hash(input.Password)
	if err != nil {
		return domain.Tokens{}, err
	}

	if user.Password != inputHashPass {
		return domain.Tokens{}, errs.ErrInvalidPass
	}

	accessToken, err := s.tokenManager.NewJWT(tokens.TokenInfo{
		UserID:   user.ID.String(),
		UserRole: domain.Role(user.UserRole),
	}, s.accessTokenTTL)
	if err != nil {
		return domain.Tokens{}, err
	}

	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return domain.Tokens{}, err
	}

	return domain.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *authService) ConfirmExpert(data domain.ConfirmExpert) error {
	err := s.authRepo.ConfirmExpert(data)
	if err != nil {
		return err
	}
	return s.ratingRepo.InsertExpertId(data.ExpertID)
}

func (s *authService) ParseAccessToken(accessToken string) (tokens.TokenInfo, error) {
	return s.tokenManager.Parse(accessToken)
}
