package common

type Pagination struct {
	Total      int
	TotalPages int
	PageSize   int
	Page       int
}

type PaginationResult[T any] struct {
	Total      int `json:",omitempty"`
	TotalPages int `json:",omitempty"`
	PageSize   int `json:",omitempty"`
	Page       int `json:",omitempty"`
	Data       []T `json:",omitempty"`
}
