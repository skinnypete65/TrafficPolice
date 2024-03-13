package converter

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/transport/rest/dto"
)

type CaseConverter struct{}

func NewCaseConverter() *CaseConverter {
	return &CaseConverter{}
}

func (c *CaseConverter) MapDtoToDomain(dto dto.Case) domain.Case {
	return domain.Case{
		Transport: domain.Transport{
			Chars:  dto.Transport.Chars,
			Num:    dto.Transport.Num,
			Region: dto.Transport.Region,
		},
		Camera: domain.Camera{
			ID: dto.Camera.ID,
		},
		Violation: domain.Violation{
			ID: dto.Violation.ID,
		},
		ViolationValue: dto.ViolationValue,
		RequiredSkill:  dto.RequiredSkill,
		Date:           dto.Date,
	}
}
