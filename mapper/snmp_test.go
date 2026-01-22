package mapper

import (
	"testing"
	"time"

	"github.com/ibm-live-project-interns/ingestor/shared/constants"
)

func TestMapSNMP_NormalCase(t *testing.T) {
	raw := []byte(`{
		"source": "router-2",
		"oid": "1.3.6.1.2.1.1.3",
		"value": "123456",
		"severity": "CRITICAL",
		"timestamp": "` + time.Now().Format(time.RFC3339) + `"
	}`)

	event, err := MapSNMP(raw)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if event.EventType != constants.EventTypeSNMP {
		t.Errorf("expected event_type=snmp, got %s", event.EventType)
	}

	if event.Severity != constants.SeverityCritical {
		t.Errorf("expected severity=critical, got %s", event.Severity)
	}

	if event.SourceHost != "router-2" {
		t.Errorf("unexpected source_host: %s", event.SourceHost)
	}
}

func TestMapSNMP_InvalidJSON(t *testing.T) {
	raw := []byte(`{this is not json}`)

	_, err := MapSNMP(raw)
	if err == nil {
		t.Fatalf("expected error for invalid JSON, got nil")
	}
}
