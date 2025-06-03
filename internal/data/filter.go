package data

import (
	"math"
	"strings"
	"time"

	"tranquara.net/internal/validator"
)

type Filter struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

type TimeFilter struct {
	StartTime    time.Time
	EndTime      time.Time
	Sort         string
	SortSafelist []string
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func (f Filter) calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

func ValidateFilter(v *validator.Validator, filter Filter) {
	v.Check(filter.Page > 0, "page", "page should be greater than 0")
	v.Check(filter.PageSize > 0, "page_size", "must be greater than zero")

	v.Check(validator.In(filter.Sort, filter.SortSafelist...), "sort", "should in the sort key list")
}

func (f Filter) sortColumn() string {
	for _, sortSafeValue := range f.SortSafelist {
		if sortSafeValue == f.Sort {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("not valid sort key " + f.Sort)
}

func (f Filter) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func (f Filter) limit() int {
	return f.PageSize
}

func (f Filter) offset() int {
	return (f.Page - 1) * f.PageSize
}
