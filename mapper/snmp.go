package mapper

import (
	"encoding/json"
	"time"
	"github.com/aishwaryagilhotra/datasource/models"
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

	return models.Event{
		EventType:      "snmp",
		SourceHost:     s.Source,
		Severity:       s.Severity,
		Message:        s.OID + " = " + s.Value,
		RawPayload:     string(rawJSON),
		EventTimestamp: ts,
	}, nil
}
