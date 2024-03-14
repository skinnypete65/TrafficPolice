package services

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/repository"
)

type DirectorService interface {
	GetCases() ([]domain.CaseStatus, error)
}

type directorService struct {
	directorRepo repository.DirectorRepo
}

func NewDirectorService(directorRepo repository.DirectorRepo) DirectorService {
	return &directorService{
		directorRepo: directorRepo,
	}
}

func (s *directorService) GetCases() ([]domain.CaseStatus, error) {
	cases, err := s.directorRepo.GetCases()
	if err != nil {
		return nil, err
	}
	if len(cases) == 0 {
		return nil, errs.ErrNoRows
	}
	return cases, nil
}
