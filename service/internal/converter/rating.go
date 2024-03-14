package converter

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/transport/rest/dto"
)

type RatingConverter struct {
}

func NewRatingConverter() *RatingConverter {
	return &RatingConverter{}
}

func (c *RatingConverter) MapDomainToDto(d domain.RatingInfo) dto.RatingInfo {
	return dto.RatingInfo{
		ExpertID:        d.ExpertID,
		Username:        d.Username,
		CompetenceSkill: d.CompetenceSkill,
		CorrectCnt:      d.CorrectCnt,
		IncorrectCnt:    d.IncorrectCnt,
	}
}

func (c *RatingConverter) MapSliceDomainToDto(domains []domain.RatingInfo) []dto.RatingInfo {
	dtos := make([]dto.RatingInfo, len(domains))
	for i := range domains {
		dtos[i] = c.MapDomainToDto(domains[i])
	}
	return dtos
}
