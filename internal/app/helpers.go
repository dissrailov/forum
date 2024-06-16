package app

import (
	"forum/internal/models"
	"net/http"
	"time"
)

func (app *Application) NewTemplateData(r *http.Request) *models.TemplateData {
	return &models.TemplateData{
		CurrentYear: time.Now().Year(),
	}
}
