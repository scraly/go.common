package request

import (
	"net/http"
	"strconv"

	"github.com/scraly/go.common/pkg/web/paginator"
)

const (
	// DefaultPerPage defines the default value for pagination
	DefaultPerPage uint = 20
)

// NewPaginatorFromRequest returns a paginator builded from an http request
func NewPaginatorFromRequest(r *http.Request) *paginator.Pagination {
	paginator := paginator.NewPaginator(1, DefaultPerPage)

	var (
		perPageRaw = r.FormValue("perPage")
		pageRaw    = r.FormValue("page")
	)

	if perPageRaw != "" {
		perPage, err := strconv.ParseUint(perPageRaw, 10, 32)
		if err != nil {
			perPage = uint64(DefaultPerPage)
		}
		paginator.PerPage = uint(perPage)
	}

	if pageRaw != "" {
		page, err := strconv.ParseUint(pageRaw, 10, 32)
		if err != nil {
			page = 1
		}
		paginator.Page = uint(page)
	}

	return paginator
}
