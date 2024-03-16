package service

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/repository"
	"TrafficPolice/internal/repository/mocks"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetCases(t *testing.T) {
	caseStatuses := []domain.CaseStatus{
		{CaseID: "case_id1", ViolationValue: "100km/h", RequiredSkill: 2, CaseDate: time.Now(), FineDecision: true},
		{CaseID: "case_id2", ViolationValue: "130km/h", RequiredSkill: 4, CaseDate: time.Now(), FineDecision: false},
		{CaseID: "case_id3", ViolationValue: "Yes", RequiredSkill: 1, CaseDate: time.Now(), FineDecision: false},
	}
	errInternal := errors.New("internal repo error")

	testCases := []struct {
		name               string
		buildDirectorRepo  func() repository.DirectorRepo
		buildCheckerRepo   func() repository.CheckerRepo
		expectedCaseStatus []domain.CaseStatus
		expectedErr        error
	}{
		{
			name: "Get statuses. Expect no error",
			buildDirectorRepo: func() repository.DirectorRepo {
				mockRepo := mocks.NewDirectorRepo(t)

				mockRepo.On("GetCases").
					Return(caseStatuses, nil)
				return mockRepo
			},
			buildCheckerRepo: func() repository.CheckerRepo {
				return mocks.NewCheckerRepo(t)
			},
			expectedCaseStatus: caseStatuses,
			expectedErr:        nil,
		},
		{
			name: "No cases. Expect ErrNoRows",
			buildDirectorRepo: func() repository.DirectorRepo {
				mockRepo := mocks.NewDirectorRepo(t)

				mockRepo.On("GetCases").
					Return(nil, nil)
				return mockRepo
			},
			buildCheckerRepo: func() repository.CheckerRepo {
				return mocks.NewCheckerRepo(t)
			},
			expectedCaseStatus: nil,
			expectedErr:        errs.ErrNoRows,
		},
		{
			name: "Unexpected error",
			buildDirectorRepo: func() repository.DirectorRepo {
				mockRepo := mocks.NewDirectorRepo(t)

				mockRepo.On("GetCases").
					Return(nil, errInternal)
				return mockRepo
			},
			buildCheckerRepo: func() repository.CheckerRepo {
				return mocks.NewCheckerRepo(t)
			},
			expectedCaseStatus: nil,
			expectedErr:        errInternal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			directorRepo := tc.buildDirectorRepo()
			checkerRepo := tc.buildCheckerRepo()

			directorService := NewDirectorService(directorRepo, checkerRepo)

			cases, err := directorService.GetCases()
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedCaseStatus, cases)
		})
	}
}
