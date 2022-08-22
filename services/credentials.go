package services

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/todanni/calendar-service/repo"
	"github.com/todanni/calendar-service/token"
)

const (
	callbackHandler = "/api/calendar/callback"
)

type CredentialsService interface {
	CredentialsCallback(w http.ResponseWriter, r *http.Request)
}

type credsService struct {
	repo   repo.CredentialsRepo
	router *mux.Router
	logger *zap.Logger
	config *oauth2.Config
}

func NewCredentialsService(repo repo.CredentialsRepo, router *mux.Router, logger *zap.Logger, config *oauth2.Config) CredentialsService {
	service := &credsService{
		repo:   repo,
		router: router,
		logger: logger,
		config: config,
	}
	service.routes()

	return service
}

func (c *credsService) routes() {
	c.router.HandleFunc(callbackHandler, c.CredentialsCallback).Methods(http.MethodPost, http.MethodGet)
}

func (c *credsService) CredentialsCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RequestURI)
	fmt.Println(r.Method)

	code := r.URL.Query().Get("code")

	tok, err := c.config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	idToken := tok.Extra("id_token").(string)
	email, err := token.ValidateToken(context.Background(), idToken)
	if err != nil {
		// TODO: log out the error and return bad response
		c.logger.Error("token validation failed: %v")
	}

	// The credentials must be persisted somehow after this
	// so they could be retrieved when requests are made
	err = c.repo.SaveUserCredentials(context.Background(), email, *tok)
	if err != nil {
		c.logger.Error("couldn't persist credentials")
	}

	w.Header().Set("Content-Type", "application/json")
	http.Redirect(w, r, "https://todanni.com/", http.StatusFound)
}
