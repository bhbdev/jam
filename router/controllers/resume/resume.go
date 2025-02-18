package resume

import (
	"net/http"

	"github.com/bhbdev/jam/lib/logger"
	"github.com/bhbdev/jam/lib/page"
	"github.com/bhbdev/jam/lib/pagination"
	"github.com/bhbdev/jam/models"
)

// const (
// 	UserPath = "/user"
// )

func ResumeList(w http.ResponseWriter, r *http.Request) {
	p := page.New(r.Context())
	p.SetPageActive("/resumes")

	params := &models.ListParams{
		Pagination: pagination.NewPagination(r),
	}

	resumes, err := p.Services.JobAppService.List(r.Context(), params)
	if err != nil {
		logger.Error("JobAppService List error", "error", err)
	}
	logger.Info("JobAppService List", "resumes", resumes)

	p.Data["Resumes"] = resumes
	p.Data["Params"] = params
	p.Render(w, "resume/list")
}
