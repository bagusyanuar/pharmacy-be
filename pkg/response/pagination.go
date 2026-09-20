package response

import (
	"math"

	"github.com/gofiber/fiber/v3"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

// Meta carries pagination metadata alongside list data.
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalData  int64 `json:"total_data"`
	TotalPages int   `json:"total_pages"`
}

// NewMeta computes total pages from total rows and limit.
func NewMeta(page, perPage int, totalData int64) *Meta {
	totalPages := 0
	if perPage > 0 {
		totalPages = int(math.Ceil(float64(totalData) / float64(perPage)))
	}
	return &Meta{
		Page:       page,
		PerPage:    perPage,
		TotalData:  totalData,
		TotalPages: totalPages,
	}
}

// PaginationParams holds clamped pagination options.
type PaginationParams struct {
	Page    int
	PerPage int
}

// Offset returns the SQL OFFSET value.
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.PerPage
}

// ParsePagination reads page and per_page query params with safe defaults and limits.
func ParsePagination(c fiber.Ctx) PaginationParams {
	page := fiber.Query(c, "page", DefaultPage)
	perPage := fiber.Query(c, "per_page", DefaultPerPage)

	if page < 1 {
		page = DefaultPage
	}
	if perPage < 1 {
		perPage = DefaultPerPage
	}
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}

	return PaginationParams{Page: page, PerPage: perPage}
}
