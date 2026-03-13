package sendlayer

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type EventsService struct {
	client *Client
}

func NewEventsService(client *Client) *EventsService {
	return &EventsService{client: client}
}

// Get retrieves events with optional filters. Pass nil to use default filters.
func (e *EventsService) Get(req *GetEventsRequest) (*EventsResponse, error) {
	if req == nil {
		req = &GetEventsRequest{}
	}
	eventOptions := map[string]bool{
		"accepted": true, "rejected": true, "delivered": true, "opened": true, "clicked": true, "unsubscribed": true, "complained": true, "failed": true,
	}
	params := map[string]string{}

	if req.StartDate != nil && req.EndDate != nil && req.EndDate.Before(*req.StartDate) {
		return nil, &SendLayerValidationError{SendLayerError{"End date must be after start date"}}
	}
	if req.StartDate != nil {
		params["StartDate"] = strconv.FormatInt(req.StartDate.Unix(), 10)
	}
	if req.EndDate != nil {
		params["EndDate"] = strconv.FormatInt(req.EndDate.Unix(), 10)
	}
	if req.Event != "" {
		if !eventOptions[req.Event] {
			return nil, &SendLayerValidationError{SendLayerError{fmt.Sprintf("Invalid event: %s", req.Event)}}
		}
		params["Event"] = req.Event
	}
	if req.MessageID != "" {
		params["MessageID"] = req.MessageID
	}
	if req.StartFrom != nil {
		params["StartFrom"] = strconv.Itoa(*req.StartFrom)
	}
	if req.RetrieveCount != nil {
		if *req.RetrieveCount <= 0 {
			return nil, &SendLayerValidationError{SendLayerError{"RetrieveCount must be greater than 0"}}
		}
		params["RetrieveCount"] = strconv.Itoa(*req.RetrieveCount)
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
