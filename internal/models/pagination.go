package models

import "math"

// PaginationParams holds pagination query parameters.
type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// DefaultPagination returns pagination params with sensible defaults.
func DefaultPagination() PaginationParams {
	return PaginationParams{
		Page:     1,
		PageSize: 10,
	}
}

// Normalize ensures pagination values are within acceptable bounds.
func (p *PaginationParams) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// Offset returns the SQL OFFSET for the current page.
func (p *PaginationParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// PaginatedResponse wraps a page of data with pagination metadata.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// NewPaginatedResponse creates a PaginatedResponse from data and total count.
func NewPaginatedResponse(data interface{}, params PaginationParams, total int64) PaginatedResponse {
	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))
	return PaginatedResponse{
		Data:       data,
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
