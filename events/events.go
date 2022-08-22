package events

import (
	"log"
	"time"

	"google.golang.org/api/calendar/v3"
)

type Event struct {
	Title     string `json:"title"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type eventsClient struct {
	calService *calendar.Service
}

func NewEventsClient(calService *calendar.Service) *eventsClient {
	return &eventsClient{
		calService: calService,
	}
}

func (c *eventsClient) RetrieveEvents() ([]Event, error) {
	evnts := make([]Event, 0)

	t := time.Now().Format(time.RFC3339)
	events, err := c.calService.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
		return nil, err
	}

	for _, item := range events.Items {
		evnts = append(evnts, Event{
			Title:     item.Summary,
			StartTime: item.Start.DateTime,
			EndTime:   item.End.DateTime,
		})
	}

	return evnts, err
}
