package data

import (
	"fmt"
	"math"
	"strings"
	"time"

	"tranquara.net/internal/validator"
)

// =============================================================================
// Configuration Constants
// =============================================================================

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// =============================================================================
// QueryFilter - Main filter struct with builder pattern
// =============================================================================

// QueryFilter provides a fluent builder pattern for constructing
// database queries with pagination, sorting, searching, and time filtering.
type QueryFilter struct {
	// Pagination
	page     int
	pageSize int

	// Sorting
	sort         string
	sortSafelist []string

	// Search
	searchQuery  string
	searchFields []string // columns to search in (for ILIKE)
	useTsVector  bool     // use PostgreSQL full-text search
	tsVectorCol  string   // name of tsvector column (e.g., "search_vector")

	// Time range
	startTime *time.Time
	endTime   *time.Time
	timeField string // column name for time filtering (e.g., "created_at")

	// Additional filters (key-value for simple equality checks)
	conditions map[string]interface{}
}

// =============================================================================
// Constructor
// =============================================================================

// NewQueryFilter creates a new QueryFilter with sensible defaults.
func NewQueryFilter() *QueryFilter {
	return &QueryFilter{
		page:       DefaultPage,
		pageSize:   DefaultPageSize,
		conditions: make(map[string]interface{}),
	}
}

// =============================================================================
// Builder Methods (Fluent API)
// =============================================================================

// WithPagination sets the page and page size for the query.
// Page must be >= 1, pageSize is clamped to MaxPageSize.
func (qf *QueryFilter) WithPagination(page, pageSize int) *QueryFilter {
	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	qf.page = page
	qf.pageSize = pageSize
	return qf
}

// WithSort sets the sort column and safelist of allowed sort columns.
// Sort can be prefixed with "-" for descending order (e.g., "-created_at").
func (qf *QueryFilter) WithSort(sort string, safelist []string) *QueryFilter {
	qf.sort = sort
	qf.sortSafelist = safelist
	return qf
}

// WithSearch enables simple ILIKE search across specified fields.
// Use this for basic substring matching.
func (qf *QueryFilter) WithSearch(query string, fields []string) *QueryFilter {
	qf.searchQuery = strings.TrimSpace(query)
	qf.searchFields = fields
	qf.useTsVector = false
	return qf
}

// WithFullTextSearch enables PostgreSQL full-text search using a tsvector column.
// This is more efficient for large datasets and supports language-aware stemming.
func (qf *QueryFilter) WithFullTextSearch(query string, tsVectorColumn string) *QueryFilter {
	qf.searchQuery = strings.TrimSpace(query)
	qf.tsVectorCol = tsVectorColumn
	qf.useTsVector = true
	return qf
}

// WithTimeRange sets a time range filter on the specified column.
// Either start or end can be nil for open-ended ranges.
func (qf *QueryFilter) WithTimeRange(start, end *time.Time, timeField string) *QueryFilter {
	qf.startTime = start
	qf.endTime = end
	qf.timeField = timeField
	return qf
}

// WithCondition adds a simple equality condition (column = value).
// Useful for filtering by foreign keys or status fields.
func (qf *QueryFilter) WithCondition(column string, value interface{}) *QueryFilter {
	qf.conditions[column] = value
	return qf
}

// =============================================================================
// Validation
// =============================================================================

// Validate checks all filter parameters and adds errors to the validator.
func (qf *QueryFilter) Validate(v *validator.Validator) {
	v.Check(qf.page > 0, "page", "must be greater than 0")
	v.Check(qf.pageSize > 0, "page_size", "must be greater than 0")
	v.Check(qf.pageSize <= MaxPageSize, "page_size", fmt.Sprintf("must not exceed %d", MaxPageSize))

	if qf.sort != "" && len(qf.sortSafelist) > 0 {
		v.Check(validator.In(qf.sort, qf.sortSafelist...), "sort", "invalid sort field")
	}

	// Validate time range if both are provided
	if qf.startTime != nil && qf.endTime != nil {
		v.Check(qf.startTime.Before(*qf.endTime) || qf.startTime.Equal(*qf.endTime),
			"time_range", "start_time must be before or equal to end_time")
	}
}

// =============================================================================
// Query Building Helpers
// =============================================================================

// SortColumn returns the column name for sorting (without the "-" prefix).
// Panics if sort is not in the safelist (should validate first).
func (qf *QueryFilter) SortColumn() string {
	if qf.sort == "" {
		return ""
	}

	for _, safe := range qf.sortSafelist {
		if safe == qf.sort {
			return strings.TrimPrefix(qf.sort, "-")
		}
	}

	panic("unsafe sort column: " + qf.sort)
}

// SortDirection returns "ASC" or "DESC" based on the sort prefix.
func (qf *QueryFilter) SortDirection() string {
	if strings.HasPrefix(qf.sort, "-") {
		return "DESC"
	}
	return "ASC"
}

// SortClause returns the complete ORDER BY clause content (e.g., "created_at DESC").
// Returns empty string if no sort is specified.
func (qf *QueryFilter) SortClause() string {
	if qf.sort == "" {
		return ""
	}
	return fmt.Sprintf("%s %s", qf.SortColumn(), qf.SortDirection())
}

// Limit returns the page size for LIMIT clause.
func (qf *QueryFilter) Limit() int {
	return qf.pageSize
}

// Offset returns the calculated offset for pagination.
func (qf *QueryFilter) Offset() int {
	return (qf.page - 1) * qf.pageSize
}

// HasSearch returns true if a search query is specified.
func (qf *QueryFilter) HasSearch() bool {
	return qf.searchQuery != ""
}

// SearchQuery returns the sanitized search query.
func (qf *QueryFilter) SearchQuery() string {
	return qf.searchQuery
}

// HasTimeRange returns true if time filtering is enabled.
func (qf *QueryFilter) HasTimeRange() bool {
	return qf.timeField != "" && (qf.startTime != nil || qf.endTime != nil)
}

// StartTime returns the start time for range filtering.
func (qf *QueryFilter) StartTime() *time.Time {
	return qf.startTime
}

// EndTime returns the end time for range filtering.
func (qf *QueryFilter) EndTime() *time.Time {
	return qf.endTime
}

// TimeField returns the column name used for time filtering.
func (qf *QueryFilter) TimeField() string {
	return qf.timeField
}

// =============================================================================
// SQL Fragment Builders
// =============================================================================

// SearchConditionSQL returns a SQL fragment for search conditions.
// Returns empty string and nil args if no search is specified.
//
// For full-text search (tsvector):
//
//	"(search_vector @@ plainto_tsquery('english', $N) OR $N = '')"
//
// For ILIKE search:
//
//	"(title ILIKE $N OR content ILIKE $N OR $N = '')"
//
// The paramIndex is the starting parameter number (e.g., 1 for $1).
func (qf *QueryFilter) SearchConditionSQL(paramIndex int) (sql string, args []interface{}) {
	if !qf.HasSearch() {
		return "", nil
	}

	if qf.useTsVector && qf.tsVectorCol != "" {
		// PostgreSQL full-text search
		sql = fmt.Sprintf("(%s @@ plainto_tsquery('english', $%d) OR $%d = '')",
			qf.tsVectorCol, paramIndex, paramIndex)
		args = []interface{}{qf.searchQuery}
		return sql, args
	}

	if len(qf.searchFields) > 0 {
		// ILIKE search across multiple fields
		var conditions []string
		searchPattern := "%" + qf.searchQuery + "%"

		for _, field := range qf.searchFields {
			conditions = append(conditions, fmt.Sprintf("%s ILIKE $%d", field, paramIndex))
		}
		// Add empty check so query works when search is empty
		conditions = append(conditions, fmt.Sprintf("$%d = ''", paramIndex))

		sql = "(" + strings.Join(conditions, " OR ") + ")"
		args = []interface{}{searchPattern}
		return sql, args
	}

	return "", nil
}

// TimeRangeConditionSQL returns a SQL fragment for time range filtering.
// Returns empty string and nil args if no time range is specified.
//
// Example output: "created_at >= $N AND created_at <= $M"
//
// The paramIndex is the starting parameter number.
func (qf *QueryFilter) TimeRangeConditionSQL(paramIndex int) (sql string, args []interface{}, nextIndex int) {
	if !qf.HasTimeRange() {
		return "", nil, paramIndex
	}

	var conditions []string
	nextIndex = paramIndex

	if qf.startTime != nil {
		conditions = append(conditions, fmt.Sprintf("%s >= $%d", qf.timeField, nextIndex))
		args = append(args, *qf.startTime)
		nextIndex++
	}

	if qf.endTime != nil {
		conditions = append(conditions, fmt.Sprintf("%s <= $%d", qf.timeField, nextIndex))
		args = append(args, *qf.endTime)
		nextIndex++
	}

	sql = strings.Join(conditions, " AND ")
	return sql, args, nextIndex
}

// FullTextRankSQL returns a SQL fragment for ordering by search relevance.
// Only applicable when using full-text search with tsvector.
func (qf *QueryFilter) FullTextRankSQL(paramIndex int) string {
	if !qf.useTsVector || qf.tsVectorCol == "" || !qf.HasSearch() {
		return ""
	}
	return fmt.Sprintf("ts_rank(%s, plainto_tsquery('english', $%d))", qf.tsVectorCol, paramIndex)
}

// =============================================================================
// Metadata Calculation
// =============================================================================

// CalculateMetadata creates pagination metadata from total record count.
func (qf *QueryFilter) CalculateMetadata(totalRecords int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  qf.page,
		PageSize:     qf.pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(qf.pageSize))),
		TotalRecords: totalRecords,
	}
}

// =============================================================================
// Getters for internal state (useful for debugging/logging)
// =============================================================================

// Page returns the current page number.
func (qf *QueryFilter) Page() int {
	return qf.page
}

// PageSize returns the page size.
func (qf *QueryFilter) PageSize() int {
	return qf.pageSize
}

// Sort returns the raw sort string.
func (qf *QueryFilter) Sort() string {
	return qf.sort
}
