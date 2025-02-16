package jobapp

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bhbdev/jam/lib/logger"
	"github.com/bhbdev/jam/lib/page"
	"github.com/bhbdev/jam/lib/pagination"
	"github.com/bhbdev/jam/models"
)

const (
	JobAppPath = "/apps"
)

func JobAppList(w http.ResponseWriter, r *http.Request) {
	p := page.New(r.Context())
	p.SetPageActive(JobAppPath)

	params := &models.ListParams{
		Pagination: pagination.NewPagination(r),
	}

	apps, err := p.Services.JobAppService.List(r.Context(), params)
	if err != nil {
		logger.Error("JobAppList error", "error", err)
	}

	p.Data["JobApps"] = apps
	p.Data["Params"] = params
	p.Render(w, "jobapp/list")
}

func JobAppDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		logger.Error("JobAppDelete invalid id", "value", r.PathValue("id"))
		http.Redirect(w, r, JobAppPath, http.StatusSeeOther)
		return
	}
	p := page.New(r.Context())
	err = p.Services.JobAppService.Delete(r.Context(), int64(id))
	if err != nil {
		logger.Error("JobAppDelete error", "error", err)
	}

	http.Redirect(w, r, JobAppPath, http.StatusSeeOther)
}

func JobAppForm(w http.ResponseWriter, r *http.Request) {
	tpl := "jobapp/form"

	appId := r.PathValue("id")
	isEditing := appId != ""

	p := page.New(r.Context())
	p.SetPageActive(JobAppPath)
	p.Data["FormAction"] = r.URL.Path
	p.Data["IsEditing"] = isEditing

	jobapp := &models.JobApp{
		DateApplied: time.Now(),
		Status:      "applied",
	}

	if isEditing {
		id, err := strconv.Atoi(appId)
		if err != nil {
			logger.Error("invalid jobapp id", "value", appId)
			http.Redirect(w, r, JobAppPath, http.StatusSeeOther)
			return
		}

		jobapp, err = p.Services.JobAppService.Get(r.Context(), id)
		if err != nil {
			logger.Error("jobapp does not exist", "id", id)
			http.Redirect(w, r, JobAppPath, http.StatusSeeOther)
			return
		}
	}
	p.Data["JobApp"] = jobapp

	if r.Method == http.MethodGet {
		p.Render(w, tpl)
		return
	}

	// handle POST
	jobapp.Position = r.PostFormValue("position")
	jobapp.Company = r.PostFormValue("company")
	jobapp.Status = r.PostFormValue("status")
	jobapp.Notes = r.PostFormValue("notes")
	dateAppliedStr := r.PostFormValue("date_applied")
	dateApplied, err := time.Parse("2006-01-02", dateAppliedStr)
	if err != nil {
		logger.Error("JobAppForm invalid date_applied", "error", err)
		p.AddError("date_applied", "Invalid date format")
	}
	jobapp.DateApplied = dateApplied

	//todo: handle file upload better.. through context? middleware?
	file, _, _ := r.FormFile("resume_file")
	if file != nil {
		err = handleResumeUpload(p, r)
		if err != nil {
			p.Render(w, tpl)
			return
		}
		jobapp.ResumeFile = p.Data["FileUpload"].(*FileUpload).FileName
	}

	// validate
	for key, err := range jobapp.Validate() {
		logger.Error("JobAppForm validation error", "key", key, "error", err)
		p.AddError(key, err)
	}

	if p.HasErrors() {
		p.Render(w, tpl)
		return
	}

	err = p.Services.JobAppService.Save(r.Context(), jobapp)
	if err != nil {
		logger.Error("JobAppForm error saving jobapp", "error", err)
		p.AddError("Database", "Error saving job application")
		p.Render(w, tpl)
		return
	}

	http.Redirect(w, r, JobAppPath, http.StatusSeeOther)
}
