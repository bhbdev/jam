package models

import (
	"slices"
	"strings"

	"gorm.io/gen/helper"
)

/*

type JobApp struct {
	ID          int64          `gorm:"primaryKey;autoIncrement:true" json:"-"` // id
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"-"`                // when created
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"-"`                // last updated
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                         // deleted at
	Status      string         `gorm:"default:active" json:"-"`                // status: applied, interview, offer, rejected, accepted
	Position    string         `json:"position"`                               // position
	Company     string         `json:"company"`                                // company
	DateApplied time.Time      `json:"date_applied"`                           // date applied
	ResumeFile  string         `json:"resume_file"`                            // resume file - this will be a link to the file
	Notes       string         `json:"notes"`                                  // notes
}
*/

func init() {
	generateModel(&ModelSpec{
		structName: "JobApp",
		tableName:  "job_apps",
		fields: []helper.Field{
			idFieldSpec,
			createdAtFieldSpec,
			deletedAtFieldSpec,
			updatedAtFieldSpec,
			&FieldSpec{
				name:    "Status",
				typ:     "string",
				gormTag: "default:applied",
				comment: "status: applied, interview, offer, rejected, accepted",
			},

			&FieldSpec{
				name:    "Position",
				typ:     "string",
				gormTag: "",
				comment: "position applied for",
			},
			&FieldSpec{
				name:    "Company",
				typ:     "string",
				gormTag: "",
				comment: "company applied to",
			},
			&FieldSpec{
				name:    "DateApplied",
				typ:     "time.Time",
				gormTag: "",
				comment: "date applied",
			},
			&FieldSpec{
				name:    "ResumeFile",
				typ:     "string",
				gormTag: "",
				comment: "resume file - this will be a link to the file under cfg.UploadsDir",
			},
			&FieldSpec{
				name:    "Notes",
				typ:     "string",
				gormTag: "",
				comment: "notes",
			},
		},
	})
}

func (m *JobApp) Validate() (errs map[string]string) {
	errs = make(map[string]string)

	if m.Position == "" {
		errs["Position"] = "Position is required"
	}
	if m.Company == "" {
		errs["Company"] = "Company is required"
	}

	JobStatus := []string{"applied", "interview", "offer", "rejected", "accepted"}
	if !slices.Contains(JobStatus, m.Status) {
		errs["Status"] = "Invalid status"
	}

	if m.ResumeFile != "" && !strings.HasSuffix(m.ResumeFile, ".pdf") {
		errs["ResumeFile"] = "Invalid file type, .pdf only"
	}

	return
}
