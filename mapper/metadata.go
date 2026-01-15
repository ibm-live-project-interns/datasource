package mapper

import (
	"encoding/json"
	"time"

	"github.com/ibm-live-project-interns/ingestor/shared/constants"
	"github.com/ibm-live-project-interns/ingestor/shared/models"
)

type MetadataInput struct {
	Entity    string `json:"entity"`
	Data      any    `json:"data"`
	Timestamp string `json:"timestamp"`
}

func MapMetadata(rawJSON []byte) (models.Event, error) {
	var m MetadataInput
	if err := json.Unmarshal(rawJSON, &m); err != nil {
		return models.Event{}, err
	}

	ts, _ := time.Parse(time.RFC3339, m.Timestamp)

	return models.Event{
		EventType:      constants.EventTypeMetadata,
		SourceHost:     m.Entity,
		SourceIP:       "0.0.0.0", // Metadata doesn't have IP
		Severity:       constants.SeverityInfo,
		Category:       "metadata",
		Message:        "Metadata update for " + m.Entity,
		RawPayload:     string(rawJSON),
		EventTimestamp: ts,
	}, nil
}
