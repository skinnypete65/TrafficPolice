package service

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/repository"
	"TrafficPolice/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestGetCase(t *testing.T) {
	defaultConsensus := 2
	userID := uuid.New()
	expertID := uuid.New()
	username := "expert"
	password := "password"
	caseID := "case_id"

	caseToSolve := domain.Case{
		ID: caseID, ViolationValue: "130km/h", RequiredSkill: 1, Date: time.Now(),
	}

	expert := domain.Expert{
		ID:              expertID.String(),
		IsConfirmed:     true,
		CompetenceSkill: 1,
		UserInfo: domain.UserInfo{
			ID:       userID,
			Username: username,
			Password: password,
			UserRole: string(domain.ExpertRole),
		},
	}

	testCases := []struct {
		name            string
		buildExpertRepo func() repository.ExpertRepo
		buildCaseRepo   func() repository.CaseRepo
		userID          uuid.UUID
		expectedCase    domain.Case
		expectedErr     error
	}{
		{
			name: "Get last not solved case",
			buildExpertRepo: func() repository.ExpertRepo {
				mockRepo := mocks.NewExpertRepo(t)
				mockRepo.On("GetExpertByUserID", userID.String()).
					Return(expert, nil)

				mockRepo.On("GetLastNotSolvedCaseID", expert.ID).
					Return(caseID, nil)

				return mockRepo
			},
			buildCaseRepo: func() repository.CaseRepo {
				mockRepo := mocks.NewCaseRepo(t)
				mockRepo.On("GetCaseByID", caseID).
					Return(caseToSolve, nil)

				return mockRepo
			},
			userID:       userID,
			expectedCase: caseToSolve,
			expectedErr:  nil,
		},
		{
			name: "Get case",
			buildExpertRepo: func() repository.ExpertRepo {
				mockRepo := mocks.NewExpertRepo(t)
				mockRepo.On("GetExpertByUserID", userID.String()).
					Return(expert, nil)

				mockRepo.On("GetLastNotSolvedCaseID", expert.ID).
					Return("", errs.ErrNoLastNotSolvedCase)

				mockRepo.On("GetNotSolvedCase", expert).
					Return(caseToSolve, nil)

				mockRepo.On("InsertNotSolvedCase", mock.Anything).
					Return(nil).
					Times(1)

				return mockRepo
			},
			buildCaseRepo: func() repository.CaseRepo {
				mockRepo := mocks.NewCaseRepo(t)
				return mockRepo
			},
			userID:       userID,
			expectedCase: caseToSolve,
			expectedErr:  nil,
		},
		{
			name: "Expert with input userID not exists",
			buildExpertRepo: func() repository.ExpertRepo {
				mockRepo := mocks.NewExpertRepo(t)
				mockRepo.On("GetExpertByUserID", userID.String()).
					Return(domain.Expert{}, errs.ErrUserNotExists)

				return mockRepo
			},
			buildCaseRepo: func() repository.CaseRepo {
				mockRepo := mocks.NewCaseRepo(t)
				return mockRepo
			},
			userID:       userID,
			expectedCase: domain.Case{},
			expectedErr:  errs.ErrUserNotExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expertRepo := tc.buildExpertRepo()
			caseRepo := tc.buildCaseRepo()

			expertService := NewExpertService(expertRepo, caseRepo, defaultConsensus)

			actualCase, err := expertService.GetCase(tc.userID.String())
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedCase, actualCase)
		})
	}
}

func TestSetCaseDecision(t *testing.T) {
	userID := uuid.New()
	expertID := uuid.New()
	username := "expert"
	password := "password"
	caseID := uuid.New()

	expert := domain.Expert{
		ID:              expertID.String(),
		IsConfirmed:     true,
		CompetenceSkill: 1,
		UserInfo: domain.UserInfo{
			ID:       userID,
			Username: username,
			Password: password,
			UserRole: string(domain.ExpertRole),
		},
	}

	positiveDecision := domain.Decision{
		CaseID:       caseID.String(),
		Expert:       expert,
		FineDecision: true,
		SolvedAt:     time.Now(),
	}
	negativeDecision := domain.Decision{
		CaseID:       caseID.String(),
		Expert:       expert,
		FineDecision: false,
	}

	testCases := []struct {
		name            string
		buildExpertRepo func() repository.ExpertRepo
		buildCaseRepo   func() repository.CaseRepo
		consensus       int
		decision        domain.Decision
		expectedInfo    domain.CaseDecisionInfo
		expectedErr     error
	}{
		{
			name: "Set decision. Not solved. Without upgrade required skill",
			buildExpertRepo: func() repository.ExpertRepo {
				mockRepo := mocks.NewExpertRepo(t)
				mockRepo.On("SetCaseDecision", mock.Anything).
					Return(nil)
				mockRepo.On("GetCaseFineDecisions", positiveDecision.CaseID).
					Return(domain.FineDecisions{PositiveDecisions: 1, NegativeDecisions: 0}, nil)
				mockRepo.On("GetExpertsCountBySkill", positiveDecision.Expert.CompetenceSkill).
					Return(4, nil)

				return mockRepo
			},
			buildCaseRepo: func() repository.CaseRepo {
				mockRepo := mocks.NewCaseRepo(t)

				return mockRepo
			},
			consensus:    2,
			decision:     positiveDecision,
			expectedInfo: domain.CaseDecisionInfo{CaseID: caseID.String(), ShouldSendFine: false, IsSolved: false},
			expectedErr:  nil,
		},
		{
			name: "Set decision. Case is solved. Fine should be sent",
			buildExpertRepo: func() repository.ExpertRepo {
				mockRepo := mocks.NewExpertRepo(t)
				mockRepo.On("SetCaseDecision", mock.Anything).
					Return(nil)
				mockRepo.On("GetCaseFineDecisions", positiveDecision.CaseID).
					Return(domain.FineDecisions{PositiveDecisions: 2, NegativeDecisions: 0}, nil)

				return mockRepo
			},
			buildCaseRepo: func() repository.CaseRepo {
				mockRepo := mocks.NewCaseRepo(t)
				mockRepo.On("SetCaseFineDecision", positiveDecision.CaseID, true, mock.Anything).
					Return(nil)

				return mockRepo
			},
			consensus:    2,
			decision:     positiveDecision,
			expectedInfo: domain.CaseDecisionInfo{CaseID: caseID.String(), ShouldSendFine: true, IsSolved: true},
			expectedErr:  nil,
		},
		{
			name: "Set decision. Case is solved. Fine should not be sent",
			buildExpertRepo: func() repository.ExpertRepo {
				mockRepo := mocks.NewExpertRepo(t)
				mockRepo.On("SetCaseDecision", mock.Anything).
					Return(nil)
				mockRepo.On("GetCaseFineDecisions", positiveDecision.CaseID).
					Return(domain.FineDecisions{PositiveDecisions: 0, NegativeDecisions: 2}, nil)

				return mockRepo
			},
			buildCaseRepo: func() repository.CaseRepo {
				mockRepo := mocks.NewCaseRepo(t)
				mockRepo.On("SetCaseFineDecision", negativeDecision.CaseID, false, mock.Anything).
					Return(nil)

				return mockRepo
			},
			consensus:    2,
			decision:     negativeDecision,
			expectedInfo: domain.CaseDecisionInfo{CaseID: caseID.String(), ShouldSendFine: false, IsSolved: true},
			expectedErr:  nil,
		},
		{
			name: "Set decision. consensus can not be reached. Upgrade required level",
			buildExpertRepo: func() repository.ExpertRepo {
				mockRepo := mocks.NewExpertRepo(t)
				mockRepo.On("SetCaseDecision", mock.Anything).
					Return(nil)
				mockRepo.On("GetCaseFineDecisions", positiveDecision.CaseID).
					Return(domain.FineDecisions{PositiveDecisions: 1, NegativeDecisions: 1}, nil)
				mockRepo.On("GetExpertsCountBySkill", positiveDecision.Expert.CompetenceSkill).
					Return(3, nil)

				return mockRepo
			},
			buildCaseRepo: func() repository.CaseRepo {
				mockRepo := mocks.NewCaseRepo(t)
				mockRepo.On("UpdateCaseRequiredSkill", positiveDecision.CaseID,
					positiveDecision.Expert.CompetenceSkill+1).
					Return(nil)
				return mockRepo
			},
			consensus:    4,
			decision:     positiveDecision,
			expectedInfo: domain.CaseDecisionInfo{CaseID: caseID.String(), ShouldSendFine: false, IsSolved: false},
			expectedErr:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expertRepo := tc.buildExpertRepo()
			caseRepo := tc.buildCaseRepo()

			expertService := NewExpertService(expertRepo, caseRepo, tc.consensus)

			actualInfo, err := expertService.SetCaseDecision(tc.decision)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedInfo, actualInfo)
		})
	}
}
