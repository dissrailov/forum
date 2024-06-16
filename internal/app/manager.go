package app

import (
	"log"
	"text/template"
)

type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	templateCache map[string]*template.Template
}

func New(infoLog, errorLog *log.Logger, templateCache map[string]*template.Template) *Application {
	return &Application{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		templateCache: templateCache,
	}
}
