package models

import (
	"bhbdev/jam/lib/pagination"
)

// Models that will be auto-migrated
var MODELS = []interface{}{
	&JobApp{},
	&User{},
}

type ListParams struct {
	Pagination *pagination.Pagination
}
