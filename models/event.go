package models

import "time"

type Event struct {
    EventType      string    `json:"event_type"`      // syslog | snmp | metadata
    SourceHost     string    `json:"source_host"`
    SourceIP       string    `json:"source_ip"`

    Severity       string    `json:"severity"`
    Category       string    `json:"category"`

    Message        string    `json:"message"`
    RawPayload     string    `json:"raw_payload"`

    EventTimestamp time.Time `json:"event_timestamp"`
}
