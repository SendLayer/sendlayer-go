package sendlayer

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type EventsService struct {
	client *Client
}

func NewEventsService(client *Client) *EventsService {
	return &EventsService{client: client}
}

func (e *EventsService) Get(
	startDate *time.Time,
	endDate *time.Time,
	event string,
	messageID string,
	startFrom *int,
	retrieveCount *int,
) (*EventsResponse, error) {
	eventOptions := map[string]bool{
		"accepted": true, "rejected": true, "delivered": true, "opened": true, "clicked": true, "unsubscribed": true, "complained": true, "failed": true,
	}
	params := map[string]string{}

	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		return nil, &SendLayerValidationError{SendLayerError{"End date must be after start date"}}
	}
	if startDate != nil {
		params["StartDate"] = strconv.FormatInt(startDate.Unix(), 10)
	}
	if endDate != nil {
		params["EndDate"] = strconv.FormatInt(endDate.Unix(), 10)
	}
	if event != "" {
		if !eventOptions[event] {
			return nil, &SendLayerValidationError{SendLayerError{fmt.Sprintf("Invalid event: %s", event)}}
		}
		params["Event"] = event
	}
	if messageID != "" {
		params["MessageID"] = messageID
	}
	if startFrom != nil {
		params["StartFrom"] = strconv.Itoa(*startFrom)
	}
	if retrieveCount != nil {
		if *retrieveCount <= 0 {
			return nil, &SendLayerValidationError{SendLayerError{"RetrieveCount must be greater than 0"}}
		}
		params["RetrieveCount"] = strconv.Itoa(*retrieveCount)
	}

	respBody, _, err := e.client.doRequest("GET", "events", nil, params)
	if err != nil {
		return nil, err
	}

	var resp EventsResponse
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
