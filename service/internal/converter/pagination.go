package converter

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/transport/rest/dto"
)

type PaginationConverter struct {
}

func NewPaginationConverter() *PaginationConverter {
	return &PaginationConverter{}
}

func (c *PaginationConverter) MapDomainToDto(pagination domain.Pagination) dto.Pagination {
	return dto.Pagination{
		Next:          pagination.Next,
		Previous:      pagination.Next,
		RecordPerPage: pagination.RecordPerPage,
		CurrentPage:   pagination.CurrentPage,
		TotalPage:     pagination.TotalPage,
	}
}
