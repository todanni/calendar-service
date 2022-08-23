package services

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	evnts "github.com/todanni/calendar-service/events"
	"github.com/todanni/calendar-service/repo"
)

const (
	eventsHandler = "/api/calendar/events"
)

type CalendarService interface {
	GetEvents(w http.ResponseWriter, r *http.Request)
}

type calendarService struct {
	repo   repo.CredentialsRepo
	router *mux.Router
	logger *zap.Logger
	config *oauth2.Config
}

func NewCalendarService(repo repo.CredentialsRepo, router *mux.Router, logger *zap.Logger, config *oauth2.Config) CalendarService {
	service := &calendarService{
		repo:   repo,
		router: router,
		logger: logger,
		config: config,
	}
	service.routes()

	return service
}

func (c *calendarService) routes() {
	c.router.HandleFunc(eventsHandler, c.GetEvents).Methods(http.MethodPost, http.MethodGet)
}

func (c *calendarService) GetEvents(w http.ResponseWriter, r *http.Request) {
	events := make([]evnts.Event, 0)
	// TODO: context should come from service
	ctx := context.Background()

	// TODO: We must validate the user's token in the future
	// and extract the email from it
	userEmail := "danni@todanni.com"
	creds, err := c.repo.GetUserCredentials(ctx, userEmail)
	if err != nil {
		// TODO: this means we have no saved credentials for this user
		// we should return bad request and redirect them to authorisation page
	}

	tokenSource := c.config.TokenSource(ctx, &creds)
	if creds.Expiry.Before(time.Now()) {
		c.logger.Info("Token is expired, renewing")
		newToken, err := tokenSource.Token()
		if err != nil {
			c.logger.Error("Couldn't get token from old token")
		}

		err = c.repo.SaveUserCredentials(ctx, userEmail, *newToken)
		if err != nil {
			c.logger.Error("Couldn't persist new token")
		}
	}

	srv, err := calendar.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		c.logger.Fatal("unable to retrieve Calendar client: %v")
	}

	eventsClient := evnts.NewEventsClient(srv)
	events, err = eventsClient.RetrieveEvents()
	if err != nil {
		c.logger.Error(err.Error())
	}

	marshalled, err := json.Marshal(events)
	if err != nil {
		c.logger.Error("couldn't marshall events")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(marshalled)
}
