package rest

import (
	"TrafficPolice/internal/converter"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/service"
	"TrafficPolice/internal/service/mocks"
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetCases(t *testing.T) {
	caseConverter := converter.NewCaseConverter()
	analyticsConverter := converter.NewAnalyticsConverter()

	caseID := uuid.New().String()
	caseStatus := domain.CaseStatus{CaseID: caseID, ViolationValue: "100km/h", RequiredSkill: 2,
		CaseDate: time.Now(), FineDecision: true}

	testCases := []struct {
		name                 string
		caseIDParam          string
		buildDirectorService func() service.DirectorService
		expectedCode         int
	}{
		{
			name:        "Get case. 200 OK",
			caseIDParam: caseID,
			buildDirectorService: func() service.DirectorService {
				mockService := mocks.NewDirectorService(t)
				mockService.On("GetCase", caseID).
					Return(caseStatus, nil)

				return mockService
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Empty case id. 400 Bad request",
			caseIDParam: "",
			buildDirectorService: func() service.DirectorService {
				mockService := mocks.NewDirectorService(t)

				return mockService
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "Case id is not uuid. 400 Bad request",
			caseIDParam: "case_id",
			buildDirectorService: func() service.DirectorService {
				mockService := mocks.NewDirectorService(t)

				return mockService
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "Case id with passed id not found. 404 Not found",
			caseIDParam: caseID,
			buildDirectorService: func() service.DirectorService {
				mockService := mocks.NewDirectorService(t)
				mockService.On("GetCase", caseID).
					Return(domain.CaseStatus{}, errs.ErrNoCase)

				return mockService
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := "/director/case"
			if tc.caseIDParam != "" {
				path += fmt.Sprintf("?id=%s", tc.caseIDParam)
			}
			handler := NewDirectorHandler(tc.buildDirectorService(), caseConverter, analyticsConverter)

			var buf bytes.Buffer

			req := httptest.NewRequest(http.MethodGet, path, &buf)
			rec := httptest.NewRecorder()

			handler.GetCase(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestExpertAnalytics(t *testing.T) {
	caseConverter := converter.NewCaseConverter()
	analyticsConverter := converter.NewAnalyticsConverter()

	expertID := uuid.New().String()

	type param struct {
		Name  string
		Param string
	}

	testCases := []struct {
		name                 string
		buildDirectorService func() service.DirectorService
		params               []param
		expectedCode         int
	}{
		{
			name: "Get analytics. Correct Params. 200 OK",
			buildDirectorService: func() service.DirectorService {
				mockService := mocks.NewDirectorService(t)
				mockService.On("GetExpertAnalytics", mock.Anything, mock.Anything, mock.Anything).
					Return([]domain.AnalyticsInterval{}, nil)

				return mockService
			},
			params: []param{
				{Name: expertIDKey, Param: expertID},
				{Name: startTimeKey, Param: "2024-01-01"},
				{Name: endTimeKey, Param: "2024-06-06"},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Start time is bad format. 400 Bad Request",
			buildDirectorService: func() service.DirectorService {
				mockService := mocks.NewDirectorService(t)

				return mockService
			},
			params: []param{
				{Name: expertIDKey, Param: expertID},
				{Name: startTimeKey, Param: "2024-01"},
				{Name: endTimeKey, Param: "2024-06-06"},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "End time is bad format. 400 Bad Request",
			buildDirectorService: func() service.DirectorService {
				mockService := mocks.NewDirectorService(t)

				return mockService
			},
			params: []param{
				{Name: expertIDKey, Param: expertID},
				{Name: startTimeKey, Param: "2024-01-01"},
				{Name: endTimeKey, Param: "2024-06"},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Expert id not passed. 400 Bad Request",
			buildDirectorService: func() service.DirectorService {
				mockService := mocks.NewDirectorService(t)

				return mockService
			},
			params: []param{
				{Name: startTimeKey, Param: "2024-01-01"},
				{Name: endTimeKey, Param: "2024-06-06"},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Expert with input id not found. 404 Not found",
			buildDirectorService: func() service.DirectorService {
				mockService := mocks.NewDirectorService(t)
				mockService.On("GetExpertAnalytics", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errs.ErrExpertNotExists)
				return mockService
			},
			params: []param{
				{Name: expertIDKey, Param: expertID},
				{Name: startTimeKey, Param: "2024-01-01"},
				{Name: endTimeKey, Param: "2024-06-06"},
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "No analytics for the expert. 204 No content",
			buildDirectorService: func() service.DirectorService {
				mockService := mocks.NewDirectorService(t)
				mockService.On("GetExpertAnalytics", mock.Anything, mock.Anything, mock.Anything).
					Return([]domain.AnalyticsInterval{}, errs.ErrNoRows)
				return mockService
			},
			params: []param{
				{Name: expertIDKey, Param: expertID},
				{Name: startTimeKey, Param: "2024-01-01"},
				{Name: endTimeKey, Param: "2024-06-06"},
			},
			expectedCode: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := "/director/analytics/expert"
			handler := NewDirectorHandler(tc.buildDirectorService(), caseConverter, analyticsConverter)

			var buf bytes.Buffer

			filteredParams := make([]param, 0)
			for i := 0; i < len(tc.params); i++ {
				if tc.params[i].Name != "" && tc.params[i].Param != "" {
					filteredParams = append(filteredParams, tc.params[i])
				}
			}

			if len(filteredParams) > 0 {
				path += "?"
				for i := 0; i < len(filteredParams)-1; i++ {
					path += fmt.Sprintf("%s=%s&", filteredParams[i].Name, filteredParams[i].Param)
				}

				last := len(filteredParams) - 1
				path += fmt.Sprintf("%s=%s", filteredParams[last].Name, filteredParams[last].Param)
			}

			req := httptest.NewRequest(http.MethodGet, path, &buf)
			rec := httptest.NewRecorder()

			handler.ExpertAnalytics(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
