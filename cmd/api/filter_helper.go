package main

import (
	"net/url"
	"time"

	"tranquara.net/internal/data"
	"tranquara.net/internal/validator"
)

// =============================================================================
// FilterOptions - Configuration for parsing filters from query strings
// =============================================================================

// FilterOptions defines the configuration for parsing a QueryFilter from URL query params.
type FilterOptions struct {
	// Pagination defaults
	DefaultPage     int
	DefaultPageSize int

	// Sorting configuration
	DefaultSort  string
	SortSafelist []string

	// Search configuration
	SearchFields   []string // For ILIKE search (e.g., []string{"title", "description"})
	TsVectorColumn string   // For full-text search (e.g., "search_vector")
	UseFullText    bool     // true = use tsvector, false = use ILIKE

	// Time range configuration
	TimeField string // Column name for time filtering (e.g., "created_at")
}

// =============================================================================
// Filter Parsing Helper
// =============================================================================

// readQueryFilter parses a QueryFilter from URL query parameters using the provided options.
// This reduces boilerplate in handlers by centralizing filter parsing logic.
//
// Query parameters supported:
//   - page: Page number (default: opts.DefaultPage or 1)
//   - page_size: Items per page (default: opts.DefaultPageSize or 20)
//   - sort: Sort field with optional "-" prefix for DESC (default: opts.DefaultSort)
//   - search: Search query string
//   - start_time: RFC3339 timestamp for range start
//   - end_time: RFC3339 timestamp for range end
func (app *application) readQueryFilter(qs url.Values, v *validator.Validator, opts FilterOptions) *data.QueryFilter {
	// Set defaults
	defaultPage := opts.DefaultPage
	if defaultPage <= 0 {
		defaultPage = data.DefaultPage
	}

	defaultPageSize := opts.DefaultPageSize
	if defaultPageSize <= 0 {
		defaultPageSize = data.DefaultPageSize
	}

	defaultSort := opts.DefaultSort
	if defaultSort == "" && len(opts.SortSafelist) > 0 {
		defaultSort = opts.SortSafelist[0]
	}

	// Build the filter
	filter := data.NewQueryFilter()

	// Pagination
	filter.WithPagination(
		app.readInt(qs, "page", defaultPage, v),
		app.readInt(qs, "page_size", defaultPageSize, v),
	)

	// Sorting
	if len(opts.SortSafelist) > 0 {
		filter.WithSort(
			app.readString(qs, "sort", defaultSort),
			opts.SortSafelist,
		)
	}

	// Search
	searchQuery := app.readString(qs, "search", "")
	if searchQuery != "" {
		if opts.UseFullText && opts.TsVectorColumn != "" {
			filter.WithFullTextSearch(searchQuery, opts.TsVectorColumn)
		} else if len(opts.SearchFields) > 0 {
			filter.WithSearch(searchQuery, opts.SearchFields)
		}
	}

	// Time range
	if opts.TimeField != "" {
		startTime := app.readTime(qs, "start_time", v)
		endTime := app.readTime(qs, "end_time", v)

		if startTime != nil || endTime != nil {
			filter.WithTimeRange(startTime, endTime, opts.TimeField)
		}
	}

	return filter
}

// readTime parses a RFC3339 timestamp from query parameters.
// Returns nil if the parameter is not provided or empty.
// Adds a validation error if the format is invalid.
func (app *application) readTime(qs url.Values, key string, v *validator.Validator) *time.Time {
	s := qs.Get(key)
	if s == "" {
		return nil
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		v.AddError(key, "must be a valid RFC3339 timestamp (e.g., 2024-01-15T09:00:00Z)")
		return nil
	}

	return &t
}

// =============================================================================
// Convenience presets for common filter configurations
// =============================================================================

// DefaultFilterOptions returns basic pagination and sorting options.
func DefaultFilterOptions(defaultSort string, sortSafelist []string) FilterOptions {
	return FilterOptions{
		DefaultPage:     1,
		DefaultPageSize: 20,
		DefaultSort:     defaultSort,
		SortSafelist:    sortSafelist,
	}
}

// SearchFilterOptions returns options with ILIKE search enabled.
func SearchFilterOptions(defaultSort string, sortSafelist, searchFields []string) FilterOptions {
	return FilterOptions{
		DefaultPage:     1,
		DefaultPageSize: 20,
		DefaultSort:     defaultSort,
		SortSafelist:    sortSafelist,
		SearchFields:    searchFields,
		UseFullText:     false,
	}
}

// FullTextFilterOptions returns options with PostgreSQL full-text search enabled.
func FullTextFilterOptions(defaultSort string, sortSafelist []string, tsVectorColumn string) FilterOptions {
	return FilterOptions{
		DefaultPage:     1,
		DefaultPageSize: 20,
		DefaultSort:     defaultSort,
		SortSafelist:    sortSafelist,
		TsVectorColumn:  tsVectorColumn,
		UseFullText:     true,
	}
}

// TimeRangeFilterOptions returns options with time range filtering.
func TimeRangeFilterOptions(defaultSort string, sortSafelist []string, timeField string) FilterOptions {
	return FilterOptions{
		DefaultPage:     1,
		DefaultPageSize: 20,
		DefaultSort:     defaultSort,
		SortSafelist:    sortSafelist,
		TimeField:       timeField,
	}
}
