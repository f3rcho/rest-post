package utils

import (
	"fmt"
	"strconv"

	"github.com/f3rcho/rest-posts/models"
)

func Pagination(pageStr, limitStr string) models.PaginationDTO {
	var err error
	var page = uint64(1)
	var limit = uint64(10)

	if pageStr != "" && limitStr != "" {
		page, err = strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			fmt.Printf("page is not valid: %v", err)
		}
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			fmt.Printf("limit is not valid: %v", err)
		}
	}

	skip := (page - 1) * limit
	return models.PaginationDTO{
		Limit: limit,
		Skip:  skip,
		Page:  page,
	}
}
