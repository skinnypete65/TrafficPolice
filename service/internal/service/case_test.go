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
)

func TestAddCase(t *testing.T) {
	caseID := uuid.New().String()
	transportID := uuid.New().String()

	testCases := []struct {
		name               string
		buildCaseRepo      func() repository.CaseRepo
		buildTransportRepo func() repository.TransportRepo
		inputCase          domain.Case
		expectedErr        error
		expectedCaseID     string
	}{
		{
			name: "Successful add case",
			buildCaseRepo: func() repository.CaseRepo {
				mockRepo := mocks.NewCaseRepo(t)

				mockRepo.On("InsertCase", mock.Anything).
					Return(caseID, nil)

				return mockRepo
			},
			buildTransportRepo: func() repository.TransportRepo {
				mockRepo := mocks.NewTransportRepo(t)

				mockRepo.On("GetTransportID", mock.Anything, mock.Anything, mock.Anything).
					Return(transportID, nil)
				return mockRepo
			},
			inputCase:      domain.Case{ID: caseID},
			expectedErr:    nil,
			expectedCaseID: caseID,
		},
		{
			name: "Transport not exists. Expect ErrNoTransport",
			buildCaseRepo: func() repository.CaseRepo {
				mockRepo := mocks.NewCaseRepo(t)
				return mockRepo
			},
			buildTransportRepo: func() repository.TransportRepo {
				mockRepo := mocks.NewTransportRepo(t)

				mockRepo.On("GetTransportID", mock.Anything, mock.Anything, mock.Anything).
					Return("", errs.ErrNoTransport)
				return mockRepo
			},
			inputCase:      domain.Case{ID: caseID},
			expectedErr:    errs.ErrNoTransport,
			expectedCaseID: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			caseRepo := tc.buildCaseRepo()
			transportRepo := tc.buildTransportRepo()

			caseService := NewCaseService(caseRepo, transportRepo)

			actualID, err := caseService.AddCase(tc.inputCase)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedCaseID, actualID)

		})
	}
}
