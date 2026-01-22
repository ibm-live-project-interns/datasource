package mapper

import (
	"testing"
	"time"

	"github.com/ibm-live-project-interns/ingestor/shared/constants"
)

func TestMapMetadata_NormalCase(t *testing.T) {
	raw := []byte(`{
		"entity": "auth-service",
		"data": {
			"version": "1.2.3"
		},
		"timestamp": "` + time.Now().Format(time.RFC3339) + `"
	}`)

	event, err := MapMetadata(raw)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if event.EventType != constants.EventTypeMetadata {
		t.Errorf("expected event_type=metadata, got %s", event.EventType)
	}

	if event.Severity != constants.SeverityInfo {
		t.Errorf("expected severity=info, got %s", event.Severity)
	}

	if event.SourceHost != "auth-service" {
		t.Errorf("unexpected source_host: %s", event.SourceHost)
	}

	if event.Category != "metadata" {
		t.Errorf("unexpected category: %s", event.Category)
	}
}

func TestMapMetadata_InvalidJSON(t *testing.T) {
	raw := []byte(`not-json`)

	_, err := MapMetadata(raw)
	if err == nil {
		t.Fatalf("expected error for invalid JSON, got nil")
	}
}
