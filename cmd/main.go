package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/todanni/template-repository/internal/config"
	"github.com/todanni/template-repository/internal/database"
	"github.com/todanni/template-repository/internal/repository"
	"github.com/todanni/template-repository/internal/server"
	"github.com/todanni/template-repository/pkg/template"
)

func main() {
	// Read config
	cfg, err := config.NewFromEnv()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Open database connection
	db, err := database.Open(cfg)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Perform any migrations needed to run the service
	err = db.AutoMigrate(&template.Template{})
	if err != nil {
		log.Error(err)
	}

	// Initialise router
	r := mux.NewRouter()

	// Create servers by passing DB connection and router
	server.NewTemplateService(repository.NewRepository(db), r)

	// Start the servers and listen
	log.Fatal(http.ListenAndServe(":8083", r))
}
