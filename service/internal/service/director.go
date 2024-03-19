package service

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/repository"
	"sort"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name DirectorService
type DirectorService interface {
	GetCase(caseID string) (domain.CaseStatus, error)
	GetExpertAnalytics(expertID string, startTime time.Time, endTime time.Time) ([]domain.AnalyticsInterval, error)
	UpdateExpertSkill(expertID string, skill int) error
}

type directorService struct {
	directorRepo repository.DirectorRepo
	checkerRepo  repository.CheckerRepo
}

func NewDirectorService(
	directorRepo repository.DirectorRepo,
	checkerRepo repository.CheckerRepo,
) DirectorService {
	return &directorService{
		directorRepo: directorRepo,
		checkerRepo:  checkerRepo,
	}
}

func (s *directorService) GetCase(caseID string) (domain.CaseStatus, error) {
	return s.directorRepo.GetCase(caseID)
}

func (s *directorService) GetExpertAnalytics(
	expertID string,
	startTime time.Time,
	endTime time.Time,
) ([]domain.AnalyticsInterval, error) {
	isExpertExists, err := s.checkerRepo.CheckExpertExists(expertID)

	if err != nil {
		return nil, err
	}
	if !isExpertExists {
		return nil, errs.ErrExpertNotExists
	}

	intervalsCases, err := s.directorRepo.GetExpertIntervalCases(expertID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	analyticsIntervals := make([]domain.AnalyticsInterval, 0)
	for date, interval := range intervalsCases {
		sort.Slice(interval, func(i, j int) bool {
			return interval[i].GotAt.Before(interval[j].GotAt)
		})

		maxConsecutive := 0
		currentConsecutive := 0
		correctCnt := 0
		incorrectCnt := 0
		unknownCnt := 0

		for i := 0; i < len(interval); i++ {
			decision := interval[i]

			if !decision.IsExpertSolve {
				unknownCnt++
				continue
			}

			if decision.ExpertFineDecision == decision.CaseFineDecision {
				correctCnt++
				currentConsecutive++
			} else if decision.ExpertFineDecision != decision.CaseFineDecision {
				incorrectCnt++
				maxConsecutive = max(maxConsecutive, currentConsecutive)
				currentConsecutive = 0
			}
		}
		maxConsecutive = max(maxConsecutive, currentConsecutive)

		analyticsIntervals = append(analyticsIntervals, domain.AnalyticsInterval{
			Date:                 date,
			AllCases:             len(interval),
			CorrectCnt:           correctCnt,
			IncorrectCnt:         incorrectCnt,
			UnknownCnt:           unknownCnt,
			MaxConsecutiveSolved: maxConsecutive,
		})
	}

	sort.Slice(analyticsIntervals, func(i, j int) bool {
		return analyticsIntervals[i].Date.Before(analyticsIntervals[j].Date)
	})

	return analyticsIntervals, nil
}

func (s *directorService) UpdateExpertSkill(expertID string, skill int) error {
	return s.directorRepo.UpdateExpertSkill(expertID, skill)
}
