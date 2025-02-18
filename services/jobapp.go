package services

import (
	"context"

	"github.com/bhbdev/jam/lib/logger"
	//"github.com/bhbdev/lib/util"
	"github.com/bhbdev/jam/models"
	"github.com/bhbdev/jam/services/repo"

	"gorm.io/gorm"
)

type JobAppService struct {
	*repo.Query
	DB *gorm.DB // this lets us use the db directly when repo.Query is not enough
}

func NewJobAppService(db *gorm.DB) *JobAppService {
	return &JobAppService{
		Query: repo.Use(db),
		DB:    db,
	}
}

func (s *JobAppService) PageServiceName() string {
	return "JobAppService"
}

func (s *JobAppService) List(ctx context.Context, params *models.ListParams) ([]*models.JobApp, error) {
	query := s.JobApp.WithContext(ctx) //.Where(s.UCType.TenantId.Eq(ctxdata.TenantId(ctx)))
	apps, err := query.
		Order(s.JobApp.CreatedAt.Desc(), s.JobApp.ID.Desc()).
		Limit(params.Pagination.PerPage()).
		Offset(params.Pagination.Offset()).
		Find()

	if err != nil {
		logger.Error("JobAppService.List error", "error", err)
		return nil, err
	}

	total, err := query.Count()
	if err != nil {
		logger.Error("JobAppService.List error counting records", "error", err)
		return nil, err
	}
	params.Pagination.SetTotal(int(total))

	return apps, nil
}

func (s *JobAppService) StatusTotals(ctx context.Context) (*models.StatusTotals, error) {

	type Totals struct {
		Total  int
		Status string
	}
	var totals []Totals

	err := s.JobApp.WithContext(ctx).
		Select(
			s.JobApp.ID.Count().As("total"),
			s.JobApp.Status.As("status"),
		).
		Group(s.JobApp.Status).
		Scan(&totals)
	if err != nil {
		logger.Error("JobAppService.StatusTotals error", "error", err)
		return nil, err
	}
	// would rather do this in the query with case/when expressions, but not sure how with gorm :\
	var statusTotals models.StatusTotals
	for _, t := range totals {
		switch t.Status {
		case "applied":
			statusTotals.Applied = t.Total
		case "interview":
			statusTotals.Interviews = t.Total
		case "offer":
			statusTotals.Offers = t.Total
		case "rejected":
			statusTotals.Rejected = t.Total
		}
	}

	total, err := s.JobApp.WithContext(ctx).Count()
	if err != nil {
		logger.Error("JobAppService.StatusTotals error counting records", "error", err)
		return nil, err
	}
	statusTotals.Total = int(total)
	logger.Info("JobAppService.StatusTotals!!!!!!!!!!!", "totals", totals)

	return &statusTotals, nil
}

func (s *JobAppService) GetLatest(ctx context.Context, limit int) ([]*models.JobApp, error) {
	apps, err := s.JobApp.WithContext(ctx).
		Order(s.JobApp.CreatedAt.Desc(), s.JobApp.ID.Desc()).
		Limit(limit).
		Find()
	if err != nil {
		logger.Error("JobAppService.GetLatest error", "error", err)
		return nil, err
	}
	return apps, nil
}

func (s *JobAppService) Get(ctx context.Context, id int) (*models.JobApp, error) {
	//app, err := s.JobApp.WithContext(ctx).Where(s.JobApp.ID.Eq(id)).First()
	app, err := s.JobApp.WithContext(ctx).GetByID(id)
	if err != nil {
		logger.Error("JobAppService.Get error", "error", err)
		return nil, err
	}
	return app, nil
}

func (s *JobAppService) Save(ctx context.Context, app *models.JobApp) error {
	//if app.ID == 0 {
	err := s.JobApp.WithContext(ctx).Save(app)
	if err != nil {
		logger.Error("JobAppService.Save error saving record", "error", err)
		return err
	}
	//}
	return nil
}

func (s *JobAppService) Delete(ctx context.Context, id int64) error {
	_, err := s.JobApp.WithContext(ctx).Where(s.JobApp.ID.Eq(id)).Delete()
	if err != nil {
		logger.Error("JobAppService.Delete error", "error", err)
		return err
	}
	return nil
}

//CREATE TABLE `job_apps` (`id` integer,`created_at` datetime,`deleted_at` datetime,`updated_at` datetime,`status` text DEFAULT "active",`position` text,`company` text,`date_applied` datetime,`resume_file` text,`notes` text,PRIMARY KEY (`id`));
// INSERT INTO `job_apps` (`position`, `company`, `date_applied`, `resume_file`, `notes`) VALUES ('Software Engineer', 'Google', '2021-01-01', 'resume.pdf', 'notes');
// INSERT INTO `job_apps` (`position`, `company`, `date_applied`, `resume_file`, `notes`) VALUES ('Software Engineer', 'Facebook', '2021-01-01', 'resume.pdf', 'notes');
