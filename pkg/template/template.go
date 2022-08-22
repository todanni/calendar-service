package template

import (
	"gorm.io/gorm"
	"net/http"
)

type Template struct {
	TemplateString string `json:"example_string"`
	gorm.Model
}

// ExampleService - service definition with all methods attached to it
type Service interface {
	// ExampleMethod description
	TemplateMethod(w http.ResponseWriter, r *http.Request)
}

// ExampleRepository - repository definition with all the methods attached to retrieving data
type Repository interface {
	//ExampleGet returns a DB record of an template
	TemplateGet() ([]Template, error)
}