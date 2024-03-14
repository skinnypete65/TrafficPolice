package services

import (
	"TrafficPolice/errs"
	"TrafficPolice/internal/config"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"log"
	"math"
	"sort"
	"time"
)

type RatingService interface {
	SetRating(caseDecision domain.CaseDecisionInfo) error
	GetRating() ([]domain.RatingInfo, error)
	RunReportPeriod(done <-chan struct{})
}

type ratingService struct {
	ratingRepo repository.RatingRepo
	ratingCfg  config.RatingConfig
}

func NewRatingService(
	ratingRepo repository.RatingRepo,
	ratingCfg config.RatingConfig,
) RatingService {
	return &ratingService{
		ratingRepo: ratingRepo,
		ratingCfg:  ratingCfg,
	}
}

func (s *ratingService) SetRating(caseDecision domain.CaseDecisionInfo) error {
	solvedDecisions, err := s.ratingRepo.GetSolvedCaseDecisions(caseDecision)
	if err != nil {
		return err
	}

	return s.ratingRepo.SetRating(solvedDecisions)
}

func (s *ratingService) GetRating() ([]domain.RatingInfo, error) {
	rating, err := s.ratingRepo.GetRating()
	if err != nil {
		return nil, err
	}
	if len(rating) == 0 {
		return nil, errs.ErrNoRows
	}

	return rating, nil
}

func (s *ratingService) RunReportPeriod(done <-chan struct{}) {
	log.Println("RunReportPeriod:", s.ratingCfg)
	ticker := time.NewTicker(s.ratingCfg.ReportPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Run report period")
			err := s.setupCompetenceSkill()
			if err != nil {
				log.Println(err)
			}
		case <-done:
			log.Println("DONE")
			return
		}
	}
}

func (s *ratingService) setupCompetenceSkill() error {
	ratings, err := s.ratingRepo.GetExpertsRating(s.ratingCfg.MinSolvedCases)
	if err != nil {
		return err
	}
	if len(ratings) < s.ratingCfg.MinExperts {
		log.Printf("Report period: Not enough experts. Min experts: %d, but got: %d",
			s.ratingCfg.MinExperts, len(ratings),
		)
		return nil
	}

	sort.Slice(ratings, func(i, j int) bool {
		return (ratings[i].CorrectCnt - ratings[i].IncorrectCnt) > (ratings[j].CorrectCnt - ratings[j].IncorrectCnt)
	})
	log.Println(ratings)
	tenPercent := int(math.Ceil(float64(len(ratings)) / 10))

	skills := make([]domain.UpdateCompetenceSkill, 0)

	// Get top 10%
	for i := 0; i < tenPercent; i++ {
		s := domain.UpdateCompetenceSkill{ExpertID: ratings[i].ExpertID, ShouldIncrease: true}
		skills = append(skills, s)
	}

	// Get last 10%
	for i := len(ratings) - 1; i > len(ratings)-1-tenPercent; i-- {
		s := domain.UpdateCompetenceSkill{ExpertID: ratings[i].ExpertID, ShouldIncrease: false}
		skills = append(skills, s)
	}

	err = s.ratingRepo.UpdateCompetenceSkills(skills)
	if err != nil {
		return err
	}
	return s.ratingRepo.ClearRating()
}
