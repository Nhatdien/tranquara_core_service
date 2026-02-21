package data

import (
	"math"
	"strings"
	"time"

	"tranquara.net/internal/validator"
)

// =============================================================================
// DEPRECATED: The Filter and TimeFilter types are deprecated.
// Use QueryFilter from query_filter.go instead, which provides:
//   - Builder pattern for fluent configuration
//   - Full-text search support (PostgreSQL tsvector)
//   - Time range filtering
//   - Better SQL generation helpers
//
// Migration guide:
//   Old: data.Filter{Page: 1, PageSize: 20, Sort: "id"}
//   New: data.NewQueryFilter().WithPagination(1, 20).WithSort("id", safelist)
// =============================================================================

// Deprecated: Use QueryFilter instead.
type Filter struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

// Deprecated: Use QueryFilter.WithTimeRange() instead.
type TimeFilter struct {
	StartTime    time.Time
	EndTime      time.Time
	Sort         string
	SortSafelist []string
}

// Metadata contains pagination metadata for list responses.
// This type is shared between the deprecated Filter and the new QueryFilter.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// Deprecated: Use QueryFilter.CalculateMetadata() instead.
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

// Deprecated: Use QueryFilter.Validate() instead.
func ValidateFilter(v *validator.Validator, filter Filter) {
	v.Check(filter.Page > 0, "page", "page should be greater than 0")
	v.Check(filter.PageSize > 0, "page_size", "must be greater than zero")

	v.Check(validator.In(filter.Sort, filter.SortSafelist...), "sort", "should in the sort key list")
}

// Deprecated: Use QueryFilter.SortColumn() instead.
func (f Filter) sortColumn() string {
	for _, sortSafeValue := range f.SortSafelist {
		if sortSafeValue == f.Sort {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("not valid sort key " + f.Sort)
}

// Deprecated: Use QueryFilter.SortDirection() instead.
func (f Filter) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

// Deprecated: Use QueryFilter.Limit() instead.
func (f Filter) limit() int {
	return f.PageSize
}

// Deprecated: Use QueryFilter.Offset() instead.
func (f Filter) offset() int {
	return (f.Page - 1) * f.PageSize
}
