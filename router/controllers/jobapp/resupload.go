package jobapp

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bhbdev/jam/lib/logger"
	"github.com/bhbdev/jam/lib/page"

	"github.com/google/uuid"
)

type FileUpload struct {
	FileName string
	FileSize int64
	FilePath string
}

func handleResumeUpload(p *page.Page, r *http.Request) (err error) {

	// Parse the multipart form, 10 MB max upload size
	r.ParseMultipartForm(10 << 20)

	// Get the file from the form
	file, handler, err := r.FormFile("resume_file")
	if err != nil {
		if err == http.ErrMissingFile {
			p.AddError("resume_file", "Please select a file to upload")
			return
		} else {
			logger.Error("Error retrieving file", "error", err)
			p.AddError("resume_file", "Error retrieving file")
		}
		return
	}
	defer file.Close()

	uuid, err := uuid.NewRandom()
	if err != nil {
		logger.Error("Error generating UUID", "error", err)
		p.AddError("resume_file", "Error generating UUID")
		return
	}

	filename := uuid.String() + filepath.Ext(handler.Filename)
	filepath := filepath.Join("assets/uploads", filename)
	if err = os.MkdirAll("assets/uploads", 0755); err != nil {
		logger.Error("Error creating directory", "error", err)
		p.AddError("resume_file", "Error saving file")
		return
	}
	logger.Info("Saving file", "filename", filename, "filepath", filepath)

	fd, err := os.Create(filepath)
	if err != nil {
		logger.Error("Error creating file", "error", err)
		p.AddError("resume_file", "Error saving file")
		return
	}
	defer fd.Close()
	if _, err = io.Copy(fd, file); err != nil {
		logger.Error("Error saving file", "error", err)
		p.AddError("resume_file", "Error saving file")
		return
	}

	p.Data["FileUpload"] = &FileUpload{
		FileName: filename,
		FileSize: handler.Size,
		FilePath: filepath,
	}

	return nil
}
