package page

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log/slog"

	"bhbdev/jam/lib/database"
	"bhbdev/jam/lib/logger"
	"bhbdev/jam/lib/middleware"
	"bhbdev/jam/services"
)

type Page struct {
	Data     map[string]interface{}
	db       database.Database
	errs     map[string]string
	tpl      *template.Template
	wrap     bool
	nav      *Nav
	Services *services.Services
}

type NavItem struct {
	Path        string
	Label       string
	IconSvgPath string
	IsActive    bool
}
type Nav struct {
	Items []NavItem
}

func (n *Nav) SetActive(path string) {
	for i := range n.Items {
		if n.Items[i].Path == path {
			n.Items[i].IsActive = true
		}
	}
}

func New(ctx context.Context) *Page {
	p := &Page{
		Data: make(map[string]interface{}),
		errs: make(map[string]string),
	}

	p.Services = ctx.Value(middleware.ContextServices).(*services.Services)
	logger.Info("page services", "services", p.Services)
	p.db = ctx.Value(middleware.ContextDatabase).(database.Database)
	p.tpl = ctx.Value(middleware.ContextTemplate).(*template.Template)
	p.wrap = true

	// if admin := ctx.Value(middleware.ContextAdminUser); admin != nil {
	// 	p.admin = admin.(*user.User)
	// }
	// p.Data["Admin"] = p.admin
	p.Data["AppName"] = "JAM"
	p.Data["AppVersion"] = "0.0.1"

	p.nav = &Nav{
		Items: []NavItem{
			{
				Path:        "/",
				Label:       "Dashboard",
				IconSvgPath: "M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6",
			},
			{
				Path:        "/apps",
				Label:       "My Apps",
				IconSvgPath: "M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01",
			},
			{
				Path:        "/resumes",
				Label:       "Resumes",
				IconSvgPath: "M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10",
			},
			{
				Path:        "/postings",
				Label:       "Job Postings",
				IconSvgPath: "M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122",
			},
		}}
	p.Data["NavItems"] = &p.nav.Items

	return p
}

func (p *Page) AddError(key, err string) {
	p.errs[key] = err
}

func (p *Page) HasErrors() bool {
	return len(p.errs) > 0
}

func (p *Page) Database() database.Database {
	return p.db
}

func (p *Page) SetPageActive(path string) {
	p.nav.SetActive(path)
}

func (p *Page) SetAutoWrap(wrap bool) {
	p.wrap = wrap
}

func (p *Page) Partial(w io.Writer, partialTemplate string) {
	// partials dont escape html!
	if err := p.tpl.ExecuteTemplate(w, partialTemplate, p.Data); err != nil {
		fmt.Fprintf(w, "Error during execution: %v", err)
		slog.Error("template execution failed", "error", err)
	}
}

func (p *Page) Render(w io.Writer, templateName string) {
	p.Data["Errors"] = p.errs
	p.Data["RenderedTemplate"] = templateName

	// reduce allocations by only using a single slice
	render := make([]string, 1, 3)
	render[0] = templateName

	// whether to wrap the desired template in the page header and footer
	if p.wrap {
		render[0] = "page/header"
		render = append(render, templateName, "page/footer")
	}

	for _, name := range render {
		if err := p.tpl.ExecuteTemplate(w, name, p.Data); err != nil {
			fmt.Fprintf(w, "Error during execution: %v", err)
			slog.Error("template execution failed", "error", err)
		}
	}
}
