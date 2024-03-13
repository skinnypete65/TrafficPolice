package converter

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/transport/rest/dto"
)

type SolvedCasesConverter struct {
}

func NewSolvedCasesConverter() *SolvedCasesConverter {
	return &SolvedCasesConverter{}
}

func (c *SolvedCasesConverter) MapParamsDtoToDomain(params dto.SolvedCasesParams) domain.SolvedCasesParams {
	return domain.SolvedCasesParams{
		CameraID:      params.CameraID,
		RequiredSkill: params.RequiredSkill,
		ViolationID:   params.ViolationID,
		StartTime:     params.StartTime,
		EndTime:       params.EndTime,
	}
}
