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

func (c *CaseConverter) MapDomainToDto(d domain.Case) dto.Case {
	return dto.Case{
		ID: d.ID,
		Transport: dto.Transport{
			ID:     d.Transport.ID,
			Chars:  d.Transport.Chars,
			Num:    d.Transport.Num,
			Region: d.Transport.Region,
			Person: &dto.Person{
				ID: d.Transport.Person.ID,
			},
		},
		Camera: dto.Camera{
			ID:           d.Camera.ID,
			CameraTypeID: d.Camera.CameraType.ID,
			Latitude:     d.Camera.Latitude,
			Longitude:    d.Camera.Longitude,
			ShortDesc:    d.Camera.ShortDesc,
		},
		Violation: dto.Violation{
			ID:         d.Violation.ID,
			Name:       d.Violation.Name,
			FineAmount: d.Violation.FineAmount,
		},
		ViolationValue: d.ViolationValue,
		RequiredSkill:  d.RequiredSkill,
		IsSolved:       d.IsSolved,
		FineDecision:   d.FineDecision,
	}
}

func (c *CaseConverter) MapCaseWithPersonToDTO(d domain.Case) dto.Case {
	return dto.Case{
		ID: d.ID,
		Transport: dto.Transport{
			ID:    d.Transport.ID,
			Chars: d.Transport.Chars,
			Num:   d.Transport.Num,
			Person: &dto.Person{
				ID:       d.Transport.Person.ID,
				PhoneNum: d.Transport.Person.PhoneNum,
				Email:    d.Transport.Person.Email,
				VkID:     d.Transport.Person.VkID,
				TgID:     d.Transport.Person.TgID,
			},
		},
		Camera: dto.Camera{
			ID:           d.Camera.ID,
			CameraTypeID: d.Camera.CameraType.ID,
			Latitude:     d.Camera.Latitude,
			Longitude:    d.Camera.Longitude,
			ShortDesc:    d.Camera.ShortDesc,
		},
		Violation: dto.Violation{
			ID:         d.Violation.ID,
			Name:       d.Violation.Name,
			FineAmount: d.Violation.FineAmount,
		},
		ViolationValue: d.ViolationValue,
		RequiredSkill:  d.RequiredSkill,
		Date:           d.Date,
		IsSolved:       d.IsSolved,
		FineDecision:   d.FineDecision,
	}
}

func (c *CaseConverter) MapCaseStatusToDto(d domain.CaseStatus) dto.CaseStatus {
	assessments := make([]dto.CaseAssessment, 0)
	for _, assessment := range d.CaseAssessments {
		assessments = append(assessments, dto.CaseAssessment{
			ExpertID: assessment.ExpertID, IsExpertSolve: assessment.IsExpertSolve, FineDecision: assessment.FineDecision},
		)
	}

	return dto.CaseStatus{
		CaseID:          d.CaseID,
		ViolationValue:  d.ViolationValue,
		RequiredSkill:   d.RequiredSkill,
		CaseDate:        d.CaseDate,
		IsSolved:        d.IsSolved,
		FineDecision:    d.FineDecision,
		CaseAssessments: assessments,
	}
}
