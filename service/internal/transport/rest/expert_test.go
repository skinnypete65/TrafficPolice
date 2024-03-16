package rest

import (
	"TrafficPolice/internal/converter"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/service"
	"TrafficPolice/internal/service/mocks"
	"TrafficPolice/internal/tokens"
	"TrafficPolice/internal/transport/rabbitmq"
	mocksmq "TrafficPolice/internal/transport/rabbitmq/mocks"
	"TrafficPolice/internal/transport/rest/middlewares"
	"TrafficPolice/pkg/imagereader"
	mocksreader "TrafficPolice/pkg/imagereader/mocks"
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCaseForExpert(t *testing.T) {
	caseConverter := converter.NewCaseConverter()
	caseDecisionConverter := converter.NewCaseDecisionConverter()
	path := "/expert/get_case"

	userID := uuid.New().String()
	tokenInfo := tokens.TokenInfo{
		UserID:   userID,
		UserRole: domain.ExpertRole,
	}

	caseForExpert := domain.Case{
		Transport: domain.Transport{Person: &domain.Person{}},
		Camera:    domain.Camera{CameraType: domain.CameraType{}},
		Violation: domain.Violation{},
	}

	testCases := []struct {
		name               string
		buildImgService    func() service.ImgService
		buildExpertService func() service.ExpertService
		buildRatingService func() service.RatingService
		buildFinePublisher func() rabbitmq.FinePublisher
		buildImageReader   func() imagereader.ImageReader
		expectedCode       int
	}{
		{
			name: "Get case for expert. 200 OK",
			buildImgService: func() service.ImgService {
				mockService := mocks.NewImgService(t)
				return mockService
			},
			buildExpertService: func() service.ExpertService {
				mockService := mocks.NewExpertService(t)
				mockService.On("GetCase", tokenInfo.UserID).
					Return(caseForExpert, nil)

				return mockService
			},
			buildRatingService: func() service.RatingService {
				mockService := mocks.NewRatingService(t)
				return mockService
			},
			buildFinePublisher: func() rabbitmq.FinePublisher {
				mockPublisher := mocksmq.NewFinePublisher(t)
				return mockPublisher
			},
			buildImageReader: func() imagereader.ImageReader {
				mockReader := mocksreader.NewImageReader(t)
				return mockReader
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Expert with input user id not found. 404 Not found",
			buildImgService: func() service.ImgService {
				mockService := mocks.NewImgService(t)
				return mockService
			},
			buildExpertService: func() service.ExpertService {
				mockService := mocks.NewExpertService(t)
				mockService.On("GetCase", tokenInfo.UserID).
					Return(caseForExpert, errs.ErrUserNotExists)

				return mockService
			},
			buildRatingService: func() service.RatingService {
				mockService := mocks.NewRatingService(t)
				return mockService
			},
			buildFinePublisher: func() rabbitmq.FinePublisher {
				mockPublisher := mocksmq.NewFinePublisher(t)
				return mockPublisher
			},
			buildImageReader: func() imagereader.ImageReader {
				mockReader := mocksreader.NewImageReader(t)
				return mockReader
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "No case for solving. 204 No content",
			buildImgService: func() service.ImgService {
				mockService := mocks.NewImgService(t)
				return mockService
			},
			buildExpertService: func() service.ExpertService {
				mockService := mocks.NewExpertService(t)
				mockService.On("GetCase", tokenInfo.UserID).
					Return(caseForExpert, errs.ErrNoCase)

				return mockService
			},
			buildRatingService: func() service.RatingService {
				mockService := mocks.NewRatingService(t)
				return mockService
			},
			buildFinePublisher: func() rabbitmq.FinePublisher {
				mockPublisher := mocksmq.NewFinePublisher(t)
				return mockPublisher
			},
			buildImageReader: func() imagereader.ImageReader {
				mockReader := mocksreader.NewImageReader(t)
				return mockReader
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name: "No case for solving. 204 No content",
			buildImgService: func() service.ImgService {
				mockService := mocks.NewImgService(t)
				return mockService
			},
			buildExpertService: func() service.ExpertService {
				mockService := mocks.NewExpertService(t)
				mockService.On("GetCase", tokenInfo.UserID).
					Return(caseForExpert, errs.ErrNoCase)

				return mockService
			},
			buildRatingService: func() service.RatingService {
				mockService := mocks.NewRatingService(t)
				return mockService
			},
			buildFinePublisher: func() rabbitmq.FinePublisher {
				mockPublisher := mocksmq.NewFinePublisher(t)
				return mockPublisher
			},
			buildImageReader: func() imagereader.ImageReader {
				mockReader := mocksreader.NewImageReader(t)
				return mockReader
			},
			expectedCode: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewExpertHandler(
				tc.buildImgService(), tc.buildExpertService(), tc.buildRatingService(),
				tc.buildFinePublisher(), tc.buildImageReader(), caseConverter, caseDecisionConverter,
			)

			var buf bytes.Buffer

			req := httptest.NewRequest(http.MethodGet, path, &buf)
			rec := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), middlewares.TokenInfoKey, tokenInfo)
			handler.GetCaseForExpert(rec, req.WithContext(ctx))
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestSetCaseDecision(t *testing.T) {
	caseConverter := converter.NewCaseConverter()
	caseDecisionConverter := converter.NewCaseDecisionConverter()
	path := "/expert/decision"

	caseID := uuid.New().String()
	userID := uuid.New()
	tokenInfo := tokens.TokenInfo{
		UserID:   userID.String(),
		UserRole: domain.ExpertRole,
	}

	expertID := uuid.New().String()
	username := "user"
	password := "pass"
	expert := domain.Expert{
		ID:              expertID,
		IsConfirmed:     true,
		CompetenceSkill: 2,
		UserInfo: domain.UserInfo{
			ID:       userID,
			Username: username,
			Password: password,
			UserRole: string(domain.ExpertRole),
		},
	}

	caseInfo := domain.Case{
		ID:        caseID,
		Transport: domain.Transport{Person: &domain.Person{}},
		Camera:    domain.Camera{CameraType: domain.CameraType{}},
		Violation: domain.Violation{},
	}

	filePath := "filepath"

	testCases := []struct {
		name               string
		buildImgService    func() service.ImgService
		buildExpertService func() service.ExpertService
		buildRatingService func() service.RatingService
		buildFinePublisher func() rabbitmq.FinePublisher
		buildImageReader   func() imagereader.ImageReader
		decision           domain.Decision
		expectedCode       int
	}{
		{
			name: "Set decision. Case is solved. Should send fine. 200 OK",
			buildImgService: func() service.ImgService {
				mockService := mocks.NewImgService(t)
				mockService.On("GetImgFilePath", mock.Anything, mock.Anything).
					Return(filePath, nil)

				return mockService
			},
			buildExpertService: func() service.ExpertService {
				mockService := mocks.NewExpertService(t)
				mockService.On("GetExpertByUserID", tokenInfo.UserID).
					Return(expert, nil)
				mockService.On("SetCaseDecision", mock.Anything).
					Return(domain.CaseDecisionInfo{CaseID: caseID, ShouldSendFine: true, IsSolved: true}, nil)
				mockService.On("GetCaseWithPersonInfo", mock.Anything).
					Return(caseInfo, nil)

				return mockService
			},
			buildRatingService: func() service.RatingService {
				mockService := mocks.NewRatingService(t)
				mockService.On("SetRating", mock.Anything).
					Return(nil)

				return mockService
			},
			buildFinePublisher: func() rabbitmq.FinePublisher {
				mockPublisher := mocksmq.NewFinePublisher(t)
				mockPublisher.On("PublishFineNotification", mock.Anything).
					Return(nil)

				return mockPublisher
			},
			buildImageReader: func() imagereader.ImageReader {
				mockReader := mocksreader.NewImageReader(t)
				mockReader.On("Read", filePath).
					Return([]byte{}, nil)
				mockReader.On("GetExtension", filePath).
					Return("jpg")

				return mockReader
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Set case decision. Case is solved. Should not send fine. 200 OK",
			buildImgService: func() service.ImgService {
				mockService := mocks.NewImgService(t)
				return mockService
			},
			buildExpertService: func() service.ExpertService {
				mockService := mocks.NewExpertService(t)
				mockService.On("GetExpertByUserID", tokenInfo.UserID).
					Return(expert, nil)
				mockService.On("SetCaseDecision", mock.Anything).
					Return(domain.CaseDecisionInfo{CaseID: caseID, ShouldSendFine: false, IsSolved: true}, nil)

				return mockService
			},
			buildRatingService: func() service.RatingService {
				mockService := mocks.NewRatingService(t)
				mockService.On("SetRating", mock.Anything).
					Return(nil)

				return mockService
			},
			buildFinePublisher: func() rabbitmq.FinePublisher {
				mockPublisher := mocksmq.NewFinePublisher(t)
				return mockPublisher
			},
			buildImageReader: func() imagereader.ImageReader {
				mockReader := mocksreader.NewImageReader(t)
				return mockReader
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Set case decision. Case is not solved. Should not send fine. 200 OK",
			buildImgService: func() service.ImgService {
				mockService := mocks.NewImgService(t)
				return mockService
			},
			buildExpertService: func() service.ExpertService {
				mockService := mocks.NewExpertService(t)
				mockService.On("GetExpertByUserID", tokenInfo.UserID).
					Return(expert, nil)
				mockService.On("SetCaseDecision", mock.Anything).
					Return(domain.CaseDecisionInfo{CaseID: caseID, ShouldSendFine: false, IsSolved: false}, nil)

				return mockService
			},
			buildRatingService: func() service.RatingService {
				mockService := mocks.NewRatingService(t)
				return mockService
			},
			buildFinePublisher: func() rabbitmq.FinePublisher {
				mockPublisher := mocksmq.NewFinePublisher(t)
				return mockPublisher
			},
			buildImageReader: func() imagereader.ImageReader {
				mockReader := mocksreader.NewImageReader(t)
				return mockReader
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Set case decision. Expert by input user id not found. 404 Not found",
			buildImgService: func() service.ImgService {
				mockService := mocks.NewImgService(t)
				return mockService
			},
			buildExpertService: func() service.ExpertService {
				mockService := mocks.NewExpertService(t)
				mockService.On("GetExpertByUserID", tokenInfo.UserID).
					Return(domain.Expert{}, errs.ErrUserNotExists)
				return mockService
			},
			buildRatingService: func() service.RatingService {
				mockService := mocks.NewRatingService(t)
				return mockService
			},
			buildFinePublisher: func() rabbitmq.FinePublisher {
				mockPublisher := mocksmq.NewFinePublisher(t)
				return mockPublisher
			},
			buildImageReader: func() imagereader.ImageReader {
				mockReader := mocksreader.NewImageReader(t)
				return mockReader
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewExpertHandler(
				tc.buildImgService(), tc.buildExpertService(), tc.buildRatingService(),
				tc.buildFinePublisher(), tc.buildImageReader(), caseConverter, caseDecisionConverter,
			)

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tc.decision)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, path, &buf)
			rec := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), middlewares.TokenInfoKey, tokenInfo)
			handler.SetCaseDecision(rec, req.WithContext(ctx))
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
