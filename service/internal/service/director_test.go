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

func mustParseTime(t *testing.T, layout string, value string) time.Time {
	parsedTime, err := time.Parse(layout, value)
	assert.NoError(t, err)
	return parsedTime
}

func TestGetExpertAnalytics(t *testing.T) {
	// DateTime = "2006-01-02 15:04:05"
	// DateOnly = "2006-01-02"

	intervalsCases := map[domain.Date][]domain.IntervalCase{
		domain.Date{Year: 2024, Month: time.January, Day: 1}: {
			{GotAt: mustParseTime(t, time.DateTime, "2024-01-01 15:00:00"), IsExpertSolve: true,
				ExpertFineDecision: true, CaseFineDecision: true},
			{GotAt: mustParseTime(t, time.DateTime, "2024-01-01 15:01:00"), IsExpertSolve: true,
				ExpertFineDecision: true, CaseFineDecision: true},
			{GotAt: mustParseTime(t, time.DateTime, "2024-01-01 15:02:00"), IsExpertSolve: true,
				ExpertFineDecision: false, CaseFineDecision: false},
			{GotAt: mustParseTime(t, time.DateTime, "2024-01-01 15:03:00"), IsExpertSolve: true,
				ExpertFineDecision: false, CaseFineDecision: true},
			{GotAt: mustParseTime(t, time.DateTime, "2024-01-01 15:05:00"), IsExpertSolve: true,
				ExpertFineDecision: true, CaseFineDecision: false},
		},
		domain.Date{Year: 2024, Month: time.January, Day: 2}: {
			{GotAt: mustParseTime(t, time.DateTime, "2024-01-02 15:00:00"), IsExpertSolve: true,
				ExpertFineDecision: true, CaseFineDecision: true},
			{GotAt: mustParseTime(t, time.DateTime, "2024-01-01 15:02:00"),
				IsExpertSolve: true, ExpertFineDecision: true, CaseFineDecision: false},
			{GotAt: mustParseTime(t, time.DateTime, "2024-01-01 15:02:00"),
				IsExpertSolve: false, ExpertFineDecision: false, CaseFineDecision: false},
		},
	}
	expertID := "expert_id"
	startTime := mustParseTime(t, time.DateOnly, "2024-01-01")
	endTime := mustParseTime(t, time.DateOnly, "2024-01-05")

	intervals := []domain.AnalyticsInterval{
		{Date: domain.Date{Year: 2024, Month: time.January, Day: 1}, AllCases: 5, CorrectCnt: 3, IncorrectCnt: 2,
			UnknownCnt: 0, MaxConsecutiveSolved: 3},
		{Date: domain.Date{Year: 2024, Month: time.January, Day: 2}, AllCases: 3, CorrectCnt: 1, IncorrectCnt: 1,
			UnknownCnt: 1, MaxConsecutiveSolved: 1},
	}

	testCases := []struct {
		name              string
		buildDirectorRepo func() repository.DirectorRepo
		buildCheckerRepo  func() repository.CheckerRepo
		expertID          string
		startTime         time.Time
		endTime           time.Time
		expectedIntervals []domain.AnalyticsInterval
		expectedErr       error
	}{
		{
			name: "Get expert intervals",
			buildDirectorRepo: func() repository.DirectorRepo {
				mockRepo := mocks.NewDirectorRepo(t)
				mockRepo.On("GetExpertIntervalCases", expertID, startTime, endTime).
					Return(intervalsCases, nil)

				return mockRepo

			},
			buildCheckerRepo: func() repository.CheckerRepo {
				mockRepo := mocks.NewCheckerRepo(t)
				mockRepo.On("CheckExpertExists", expertID).
					Return(true, nil)

				return mockRepo
			},
			expertID:          expertID,
			startTime:         startTime,
			endTime:           endTime,
			expectedIntervals: intervals,
			expectedErr:       nil,
		},
		{
			name: "Expert not exists. Expect ErrExpertNotExists",
			buildDirectorRepo: func() repository.DirectorRepo {
				mockRepo := mocks.NewDirectorRepo(t)
				return mockRepo

			},
			buildCheckerRepo: func() repository.CheckerRepo {
				mockRepo := mocks.NewCheckerRepo(t)
				mockRepo.On("CheckExpertExists", expertID).
					Return(false, nil)

				return mockRepo
			},
			expertID:          expertID,
			startTime:         startTime,
			endTime:           endTime,
			expectedIntervals: nil,
			expectedErr:       errs.ErrExpertNotExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			directorRepo := tc.buildDirectorRepo()
			checkerRepo := tc.buildCheckerRepo()

			directorService := NewDirectorService(directorRepo, checkerRepo)

			actualIntervals, err := directorService.GetExpertAnalytics(tc.expertID, tc.startTime, tc.endTime)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedIntervals, actualIntervals)
		})
	}
}
