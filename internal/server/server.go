package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/todanni/template-repository/pkg/template"
)

func NewTemplateService(repo template.Repository, router *mux.Router) template.Service {
	server := &templateService{
		repo:   repo,
		router: router,
	}
	server.routes()
	return server
}

type templateService struct {
	repo   template.Repository
	router *mux.Router
}

func (s *templateService) TemplateMethod(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
