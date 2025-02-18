package controllers

import (
	"net/http"

	"github.com/bhbdev/jam/lib/page"
	"github.com/bhbdev/jam/models"
)

func Index(w http.ResponseWriter, r *http.Request) {
	p := page.New(r.Context())
	p.SetPageActive(r.URL.Path)

	// status totals
	p.Data["StatusTotals"] = &models.StatusTotals{}
	statusTotals, _ := p.Services.JobAppService.StatusTotals(r.Context())
	if statusTotals != nil {
		p.Data["StatusTotals"] = statusTotals
	}

	// latest job apps
	p.Data["LatestApps"] = []*models.JobApp{}
	latestApps, _ := p.Services.JobAppService.GetLatest(r.Context(), 5)
	if latestApps != nil {
		p.Data["LatestApps"] = latestApps
	}

	p.Render(w, "index")
}
