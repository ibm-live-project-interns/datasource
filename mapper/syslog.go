package mapper

import (
	"encoding/json"
	"time"
	"github.com/aishwaryagilhotra/datasource/models"
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

	return models.Event{
		EventType:      "syslog",
		SourceHost:     s.Host,
		Severity:       s.Severity,
		Message:        s.Message,
		RawPayload:     string(rawJSON),
		EventTimestamp: ts,
	}, nil
}
