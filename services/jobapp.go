package services

import (
	"context"

	"bhbdev/jam/lib/logger"
	//"bhbdev/lib/util"
	"bhbdev/jam/models"
	"bhbdev/jam/services/repo"

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
