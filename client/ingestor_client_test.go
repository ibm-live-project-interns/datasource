package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ibm-live-project-interns/ingestor/shared/models"
)

func TestSendEvent_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewIngestorClient(server.URL)

	event := models.Event{
		EventType:      "syslog",
		SourceHost:     "router-1",
		SourceIP:       "192.168.1.1",
		Severity:       "critical",
		Category:       "network",
		Message:        "Interface down",
		EventTimestamp: time.Now(),
	}

	err := client.SendEvent(event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSendEvent_ClientError_NoRetry(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewIngestorClient(server.URL)

	event := models.Event{
		EventType:      "syslog",
		SourceHost:     "router-1",
		SourceIP:       "192.168.1.1",
		Severity:       "critical",
		Category:       "network",
		Message:        "Bad event",
		EventTimestamp: time.Now(),
	}

	_ = client.SendEvent(event)

	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}

func TestHealthCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewIngestorClient(server.URL)

	err := client.HealthCheck()
	if err != nil {
		t.Fatalf("expected health check to succeed, got %v", err)
	}
}
