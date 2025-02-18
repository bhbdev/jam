package resume

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bhbdev/jam/lib/logger"
	"github.com/bhbdev/jam/lib/page"

	"github.com/google/uuid"
)

type FileUpload struct {
	FileName string
	FileSize int64
	FilePath string
}

func UploadDelete(w http.ResponseWriter, r *http.Request) {
	fileId := r.PathValue("id")
	if !strings.HasSuffix(fileId, ".pdf") {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}
	p := page.New(r.Context())

	if err := os.Remove(filepath.Join("assets/uploads", fileId)); err != nil {
		logger.Error("Error deleting file", "error", err)
	}
	p.Data["Swap"] = true
	p.Data["FileUpload"] = &FileUpload{}
	p.Render(w, "resume/upload")
	return
}

// func JobAppForm(w http.ResponseWriter, r *http.Request) {
func Upload(w http.ResponseWriter, r *http.Request) {

	tpl := "resume/upload"
	errs := make(map[string]string)

	p := page.New(r.Context())
	p.SetAutoWrap(false)
	p.Data["FileUpload"] = &FileUpload{}

	if r.Method == http.MethodGet {
		p.Render(w, tpl)
		return
	}

	// Parse the multipart form, 10 MB max upload size
	r.ParseMultipartForm(10 << 20)

	// Get the file from the form
	file, handler, err := r.FormFile("upload")
	if err != nil {
		if err == http.ErrMissingFile {
			errs["resume_file"] = "Please select a file to upload"
		} else {
			logger.Error("Error retrieving file", "error", err)
			errs["Upload"] = "Error saving file"
		}
	}
	defer file.Close()

	uuid, err := uuid.NewRandom()
	if err != nil {
		logger.Error("Error generating UUID", "error", err)
		errs["Upload"] = "Error saving file"
	}

	filename := uuid.String() + filepath.Ext(handler.Filename)
	filepath := filepath.Join("assets/uploads", filename)
	if err = os.MkdirAll("assets/uploads", 0755); err != nil {
		logger.Error("Error creating directory", "error", err)
		errs["Upload"] = "Error saving file"
	}

	fd, err := os.Create(filepath)
	if err != nil {
		logger.Error("Error creating file", "error", err)
		errs["Upload"] = "Error saving file"
	}
	defer fd.Close()
	if _, err = io.Copy(fd, file); err != nil {
		logger.Error("Error saving file", "error", err)
		errs["Upload"] = "Error saving file"
	}

	for key, err := range errs {
		p.AddError(key, err)
	}
	if p.HasErrors() {
		p.Render(w, tpl)
		return
	}

	p.Data["Swap"] = true
	p.Data["FileUpload"] = &FileUpload{
		FileName: filename,
		FileSize: handler.Size,
		FilePath: filepath,
	}

	p.Render(w, tpl)
}
