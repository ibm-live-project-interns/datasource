package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ibm-live-project-interns/ingestor/shared/models"
)

// IngestorClient handles communication with Ingestor Core
type IngestorClient struct {
	baseURL    string
	httpClient *http.Client
	maxRetries int
	retryDelay time.Duration
}

// NewIngestorClient creates a new client for Ingestor Core
func NewIngestorClient(baseURL string) *IngestorClient {
	return &IngestorClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		maxRetries: 3,
		retryDelay: 2 * time.Second,
	}
}

// SendEvent sends an event to Ingestor Core with retry logic
func (c *IngestorClient) SendEvent(event models.Event) error {
	// Validate event before sending
	if err := event.Validate(); err != nil {
		return fmt.Errorf("event validation failed: %w", err)
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	url := fmt.Sprintf("%s/ingest/event", c.baseURL)

	var lastErr error
	for attempt := 1; attempt <= c.maxRetries; attempt++ {
		resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: request failed: %w", attempt, err)
			if attempt < c.maxRetries {
				time.Sleep(c.retryDelay * time.Duration(attempt)) // Exponential backoff
				continue
			}
			break
		}
		defer resp.Body.Close()

		// Read response body
		bodyBytes, _ := io.ReadAll(resp.Body)

		// Check if request was successful
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}

		// Handle different error status codes
		lastErr = fmt.Errorf("attempt %d: ingestor returned status %d: %s",
			attempt, resp.StatusCode, string(bodyBytes))

		// Don't retry on client errors (4xx), only on server errors (5xx)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			break
		}

		if attempt < c.maxRetries {
			time.Sleep(c.retryDelay * time.Duration(attempt))
		}
	}

	return fmt.Errorf("failed to send event after %d attempts: %w", c.maxRetries, lastErr)
}

// SendEvents sends multiple events in batch
func (c *IngestorClient) SendEvents(events []models.Event) []error {
	errors := make([]error, 0)

	for i, event := range events {
		if err := c.SendEvent(event); err != nil {
			errors = append(errors, fmt.Errorf("event %d failed: %w", i, err))
		}
	}

	return errors
}

// HealthCheck checks if Ingestor Core is reachable
func (c *IngestorClient) HealthCheck() error {
	url := fmt.Sprintf("%s/health", c.baseURL)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}

	return nil
}
