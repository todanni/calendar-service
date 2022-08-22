package services

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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
	creds, err := c.repo.GetUserCredentials(context.Background(), userEmail)
	if err != nil {
		// TODO: this means we have no saved credentials for this user
		// we should return bad request and redirect them to authorisation page
	}
	client := c.config.Client(ctx, &creds)

	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	//calendar.NewService(context.Background(), option.WithTokenSource(c.repo))

	c.config.TokenSource(ctx, &creds)

	if err != nil {
		log.Fatalf("unable to retrieve Calendar client: %v", err)
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

	w.Write(marshalled)
}
