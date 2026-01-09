package mapper

import (
	"encoding/json"
	"time"
	"github.com/aishwaryagilhotra/datasource/models"
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
		EventType:      "metadata",
		SourceHost:     m.Entity,
		Message:        "Metadata update",
		RawPayload:     string(rawJSON),
		EventTimestamp: ts,
	}, nil
}
