package clients

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type EventsClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewEventsClient(baseURL string) *EventsClient {
	return &EventsClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *EventsClient) ValidateEvent(ctx context.Context, eventID int64, authHeader string) error {
	if c.baseURL == "" {
		return fmt.Errorf("API_BASE_URL is not configured")
	}
	if eventID == 0 {
		return fmt.Errorf("eventID is required")
	}

	url := fmt.Sprintf("%s/events/%d", c.baseURL, eventID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("event not found")
	case http.StatusUnauthorized, http.StatusForbidden:
		return fmt.Errorf("unauthorized to validate event")
	default:
		return fmt.Errorf("event validation failed with status %d", resp.StatusCode)
	}
}
