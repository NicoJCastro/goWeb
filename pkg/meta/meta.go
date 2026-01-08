package meta

import (
	"os"
	"strconv"
)

//Meta de request. Manejamos la informacion del body, ej: cantidad de registros, pagina actual, etc.

type Meta struct {
	TotalCount int `json:"total_count"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	PageCount  int `json:"page_count"`
}

func New(page, perPage, totalCount int) (*Meta, error) {

	if perPage <= 0 {
		var err error
		perPage, err = strconv.Atoi(os.Getenv("PAGINATION_LIMIT"))
		if err != nil {
			return nil, err
		}
	}

	pageCount := 0
	if totalCount >= 0 {
		pageCount = (totalCount + perPage - 1) / perPage
		if page > pageCount {
			page = pageCount
		}
	}

	if page < 1 {
		page = 1
	}

	return &Meta{
		Page:       page,
		PerPage:    perPage,
		TotalCount: totalCount,
		PageCount:  pageCount,
	}, nil

}

func (p *Meta) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func (p *Meta) Limit() int {
	return p.PerPage
}
