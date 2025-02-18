package models

import (
	"github.com/bhbdev/jam/lib/pagination"
)

// Models that will be auto-migrated
var MODELS = []interface{}{
	&JobApp{},
	&User{},
}

type ListParams struct {
	Pagination *pagination.Pagination
}

type StatusTotals struct {
	Total      int
	Applied    int
	Interviews int
	Offers     int
	Rejected   int
}
