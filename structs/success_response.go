package structs

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// pagination meta data
type PaginationMeta struct {
	TotalData   int64 `json:"total_data"`
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	Limit       int   `json:"limit"`
}

// Response struct for paginated response, used to return data with pagination meta data in API responses
type PaginatedResponse struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
	Data       any            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}