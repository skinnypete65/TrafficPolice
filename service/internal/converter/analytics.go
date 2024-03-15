package converter

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/transport/rest/dto"
	"fmt"
)

type AnalyticsConverter struct {
}

func NewAnalyticsConverter() *AnalyticsConverter {
	return &AnalyticsConverter{}
}

func (c *AnalyticsConverter) MapDomainToDto(d domain.AnalyticsInterval) dto.AnalyticsInterval {
	return dto.AnalyticsInterval{
		Date:                 fmt.Sprintf("%04d-%02d-%02d", d.Date.Year, d.Date.Month, d.Date.Day),
		AllCases:             d.AllCases,
		CorrectCnt:           d.CorrectCnt,
		IncorrectCnt:         d.IncorrectCnt,
		UnknownCnt:           d.UnknownCnt,
		MaxConsecutiveSolved: d.MaxConsecutiveSolved,
	}
}

func (c *AnalyticsConverter) MapDomainsToDtos(d []domain.AnalyticsInterval) []dto.AnalyticsInterval {
	dtos := make([]dto.AnalyticsInterval, len(d))
	for i := range d {
		dtos[i] = c.MapDomainToDto(d[i])
	}
	return dtos
}
