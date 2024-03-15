package app

import "TrafficPolice/internal/converter"

type converters struct {
	auth          *converter.AuthConverter
	camera        *converter.CameraConverter
	caseConverter *converter.CaseConverter
	caseDecision  *converter.CaseDecisionConverter
	pagination    *converter.PaginationConverter
	solvedCases   *converter.SolvedCasesConverter
	userInfo      *converter.UserInfoConverter
	analytics     *converter.AnalyticsConverter
	rating        *converter.RatingConverter
}

func newConverters() *converters {
	return &converters{
		auth:          converter.NewAuthConverter(),
		camera:        converter.NewCameraConverter(),
		caseConverter: converter.NewCaseConverter(),
		caseDecision:  converter.NewCaseDecisionConverter(),
		pagination:    converter.NewPaginationConverter(),
		solvedCases:   converter.NewSolvedCasesConverter(),
		userInfo:      converter.NewUserInfoConverter(),
		analytics:     converter.NewAnalyticsConverter(),
		rating:        converter.NewRatingConverter(),
	}
}
