package mapper

import (
	"encoding/json"
	"time"

	"github.com/ibm-live-project-interns/ingestor/shared/constants"
	"github.com/ibm-live-project-interns/ingestor/shared/models"
)

type SyslogInput struct {
	Host      string `json:"host"`
	Severity  string `json:"severity"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func MapSyslog(rawJSON []byte) (models.Event, error) {
	var s SyslogInput
	if err := json.Unmarshal(rawJSON, &s); err != nil {
		return models.Event{}, err
	}

	ts, _ := time.Parse(time.RFC3339, s.Timestamp)

	// Normalize severity to standard format
	severity := normalizeSeverity(s.Severity)

	// Resolve the hostname to an IP address using the resolver
	sourceIP := ResolveHostIP(s.Host)

	return models.Event{
		EventType:      constants.EventTypeSyslog,
		SourceHost:     s.Host,
		SourceIP:       sourceIP,
		Severity:       severity,
		Category:       "system",
		Message:        s.Message,
		RawPayload:     string(rawJSON),
		EventTimestamp: ts,
	}, nil
}

// normalizeSeverity converts various severity formats to standard format
func normalizeSeverity(severity string) string {
	switch severity {
	case "ERROR", "CRITICAL", "ALERT", "EMERGENCY":
		return constants.SeverityCritical
	case "WARN", "WARNING":
		return constants.SeverityHigh
	case "NOTICE":
		return constants.SeverityMedium
	case "DEBUG":
		return constants.SeverityLow
	case "INFO", "INFORMATIONAL":
		return constants.SeverityInfo
	default:
		return constants.SeverityInfo
	}
}
