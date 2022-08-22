package server

import "net/http"

func (s *templateService) routes() {
	// GET an template
	s.router.HandleFunc("/api/template", s.TemplateMethod).Methods(http.MethodGet)
}