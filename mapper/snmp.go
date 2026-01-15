package mapper

import (
	"encoding/json"
	"time"

	"github.com/ibm-live-project-interns/ingestor/shared/constants"
	"github.com/ibm-live-project-interns/ingestor/shared/models"
)

type SNMPInput struct {
	Source    string `json:"source"`
	OID       string `json:"oid"`
	Value     string `json:"value"`
	Severity  string `json:"severity"`
	Timestamp string `json:"timestamp"`
}

func MapSNMP(rawJSON []byte) (models.Event, error) {
	var s SNMPInput
	if err := json.Unmarshal(rawJSON, &s); err != nil {
		return models.Event{}, err
	}

	ts, _ := time.Parse(time.RFC3339, s.Timestamp)

	// Normalize severity to standard format
	severity := normalizeSeverity(s.Severity)

	return models.Event{
		EventType:      constants.EventTypeSNMP,
		SourceHost:     s.Source,
		SourceIP:       "0.0.0.0", // TODO: Resolve from source
		Severity:       severity,
		Category:       "network",
		Message:        s.OID + " = " + s.Value,
		RawPayload:     string(rawJSON),
		EventTimestamp: ts,
	}, nil
}
