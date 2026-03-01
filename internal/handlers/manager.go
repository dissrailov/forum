package handlers

import (
	"forum/internal/ai"
	"forum/internal/app"
	"forum/internal/service"
)

type HandlerApp struct {
	service   service.ServiceI
	aiService *ai.Service
	*app.Application
}

func New(s service.ServiceI, a *app.Application, aiSvc *ai.Service) *HandlerApp {
	return &HandlerApp{
		service:     s,
		aiService:   aiSvc,
		Application: a,
	}
}
