package handle

import (
	"embed"
	"html/template"
	"log/slog"
)

var (
	//go:embed ./templates
	templateFS embed.FS
	templates  = template.Must(template.ParseFS(templateFS, "templates/*.html"))
)

const (
	queryAuthRequestID = "authRequestID"
)

func errMsg(err error) string {
	if err == nil {
		return ""
	}
	errMsg := err.Error()
	slog.Error(errMsg)
	return errMsg
}
