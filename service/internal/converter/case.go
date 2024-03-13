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

func (c *CaseConverter) MapDomainsToDto(cases []domain.Case) []dto.Case {
	dtos := make([]dto.Case, 0, len(cases))
	for _, c := range cases {
		cDto := dto.Case{
			ID: c.ID,
			Transport: dto.Transport{
				ID:     c.Transport.ID,
				Chars:  c.Transport.Chars,
				Num:    c.Transport.Num,
				Region: c.Transport.Region,
			},
			Camera: dto.Camera{
				ID:           c.Camera.ID,
				CameraTypeID: c.Camera.CameraType.ID,
				Latitude:     c.Camera.Latitude,
				Longitude:    c.Camera.Longitude,
				ShortDesc:    c.Camera.ShortDesc,
			},
			Violation: dto.Violation{
				ID:         c.Violation.ID,
				Name:       c.Violation.Name,
				FineAmount: c.Violation.FineAmount,
			},
			ViolationValue: c.ViolationValue,
			RequiredSkill:  c.RequiredSkill,
			IsSolved:       c.IsSolved,
			FineDecision:   c.FineDecision,
			Date:           c.Date,
		}

		dtos = append(dtos, cDto)
	}

	return dtos
}
