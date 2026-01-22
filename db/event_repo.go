// NOTE:
// This DB package is optional and NOT part of the default datasource runtime.
// Persistence is handled downstream by Ingestor / API Gateway.
// This package is retained for local debugging or future extensions.package db

import (
	"database/sql"

	"github.com/ibm-live-project-interns/ingestor/shared/models"
)

func InsertEvent(db *sql.DB, e models.Event) error {
	query := `
		INSERT INTO events (
			event_type,
			source_host,
			source_ip,
			severity,
			category,
			message,
			raw_payload,
			event_timestamp
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`

	var sourceIP interface{}
	if e.SourceIP == "" {
		sourceIP = nil
	} else {
		sourceIP = e.SourceIP
	}

	_, err := db.Exec(
		query,
		e.EventType,
		e.SourceHost,
		sourceIP, // ✅ NULL-safe
		e.Severity,
		e.Category,
		e.Message,
		e.RawPayload,
		e.EventTimestamp,
	)

	return err
}
