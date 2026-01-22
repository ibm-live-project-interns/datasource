package mapper

import (
	"testing"
	"time"

	"github.com/ibm-live-project-interns/ingestor/shared/constants"
)

func TestMapSyslog_NormalCase(t *testing.T) {
	raw := []byte(`{
		"host": "router-1",
		"severity": "ERROR",
		"message": "Interface Gi0/1 down",
		"timestamp": "` + time.Now().Format(time.RFC3339) + `"
	}`)

	event, err := MapSyslog(raw)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if event.EventType != constants.EventTypeSyslog {
		t.Errorf("expected event_type=syslog, got %s", event.EventType)
	}

	if event.Severity != constants.SeverityCritical {
		t.Errorf("expected severity=critical, got %s", event.Severity)
	}

	if event.SourceHost != "router-1" {
		t.Errorf("unexpected source_host: %s", event.SourceHost)
	}
}

func TestMapSyslog_InvalidJSON(t *testing.T) {
	raw := []byte(`{invalid json}`)

	_, err := MapSyslog(raw)
	if err == nil {
		t.Fatalf("expected error for invalid JSON, got nil")
	}
}
