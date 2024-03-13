package converter

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/transport/rest/dto"
)

type AuthConverter struct {
}

func NewAuthConverter() *AuthConverter {
	return &AuthConverter{}
}

func (c *AuthConverter) MapConfirmExpertDtoToDomain(confirm dto.ConfirmExpertInput) domain.ConfirmExpert {
	return domain.ConfirmExpert{
		ExpertID:    confirm.ExpertID,
		IsConfirmed: confirm.IsConfirmed,
	}
}
