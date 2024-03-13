package converter

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/transport/rest/dto"
)

type CaseDecisionConverter struct {
}

func NewCaseDecisionConverter() *CaseDecisionConverter {
	return &CaseDecisionConverter{}
}

func (c *CaseDecisionConverter) MapDtoToDomain(decision dto.Decision, expert domain.Expert) domain.Decision {
	return domain.Decision{
		CaseID:       decision.CaseID,
		Expert:       expert,
		FineDecision: decision.FineDecision,
	}
}
