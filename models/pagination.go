package models

type PaginationDTO struct {
	Limit uint64
	Skip  uint64
	Page  uint64
}
